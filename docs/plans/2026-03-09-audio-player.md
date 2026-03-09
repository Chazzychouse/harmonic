# Audio Player

## Overview

Add audio playback to Harmonic. The player uses the HTML5 `<audio>` element on the frontend for playback, a Svelte store for state, and a thin Go HTTP handler to stream files from disk into the webview. No external libraries required.

```
┌─────────────────────────────────────────┐
│  Toolbar (52px)                         │
├─────────────────────────────────────────┤
│                                         │
│  Library / TrackList (scrollable)       │
│                                         │
├─────────────────────────────────────────┤
│  NowPlaying bar (64px, fixed bottom)    │
│  [track info] [◀ ▶❚❚ ▶▶] [seek] [vol] │
└─────────────────────────────────────────┘
```

The NowPlaying bar only appears when a track is loaded.

---

## Why no libraries

| Concern | Why built-in is enough |
|---------|----------------------|
| Playback | `<audio>` handles mp3, flac, wav, ogg, opus, m4a, aac natively in WebKit (Linux) and Chromium (Windows). |
| Seeking | `http.ServeFile` in Go handles HTTP range requests, so the `<audio>` element can seek without buffering the entire file. |
| State | A single Svelte writable store. The state shape is small (track, queue, index, paused, time, volume). |
| UI | Transport controls are SVG icons + a styled range input. The existing design system covers everything. |
| Visualization | Not in scope. If added later, wavesurfer.js is the obvious choice but is a separate concern. |

Something like Howler.js adds a Web Audio API abstraction layer we don't need. The `<audio>` element is simpler and sufficient.

---

## Step 1: Go — Audio file streaming handler

### Problem

The webview cannot open `file://` paths directly. The `<audio>` element needs an HTTP URL to play from.

### Solution

Add a custom `http.Handler` to the Wails asset server. Requests that don't match an embedded frontend asset fall through to this handler.

### New file: `internal/fs/handler.go`

Create a function that returns an `http.Handler`:

```go
func NewAudioHandler() http.Handler
```

The handler:

1. Only accepts `GET` requests. Return `405 Method Not Allowed` otherwise.
2. Reads the `path` query parameter: e.g. `/audio?path=/home/user/music/song.mp3`.
3. Validates the file extension against the existing `audioExtensions` map from `dir.go`. Return `400 Bad Request` if the extension isn't recognized.
4. Checks that the path exists and is a regular file (not a directory, not a symlink to something weird). Return `404 Not Found` if not.
5. Calls `http.ServeFile(w, r, filePath)`. This handles:
   - `Content-Type` based on file extension
   - `Content-Length`
   - `Accept-Ranges: bytes` and range request responses (HTTP 206) — required for seeking
   - Conditional requests (`If-Modified-Since`, `ETag`)

That's the entire handler. No buffering, no goroutines, no streaming logic — `http.ServeFile` does it all.

### Security note

The handler should reject paths that don't have a recognized audio extension. This prevents the frontend from being used to exfiltrate arbitrary files. No further sandboxing is needed for a local-only desktop app.

### Modify: `main.go`

Wire the handler into the existing `AssetServer` config. The `assetserver.Options` struct has a `Handler` field:

```go
AssetServer: &assetserver.Options{
    Assets:  assets,
    Handler: fs.NewAudioHandler(),
},
```

Wails routes: embedded frontend assets are served first; any request that doesn't match falls through to `Handler`. Since the handler responds to `/audio?path=...` and the frontend assets live at `/`, there's no conflict.

### How to verify

After implementing, you can test without any frontend changes:

1. Run `wails dev`
2. Open the webview devtools (right-click → Inspect)
3. In the console: `new Audio("/audio?path=/absolute/path/to/some.mp3").play()`
4. You should hear audio. Check the Network tab to confirm range requests (HTTP 206) work — scrub the audio by setting `.currentTime` on the element.

---

## Step 2: Frontend — Player store

### New file: `frontend/src/lib/stores/player.ts`

A Svelte writable store that holds all playback state, plus exported functions to mutate it.

### State shape

```typescript
interface AudioFile {
  title: string;
  file_path: string;
  ext: string;
}

interface PlayerState {
  track: AudioFile | null;  // currently loaded track
  queue: AudioFile[];       // ordered playlist (snapshot of filteredTracks at time of play)
  queueIndex: number;       // current position in queue (-1 if nothing loaded)
  paused: boolean;          // true = paused, false = playing
  currentTime: number;      // seconds, updated by <audio> timeupdate events
  duration: number;         // seconds, set on loadedmetadata
  volume: number;           // 0.0 – 1.0
}
```

### Initial state

```typescript
const initial: PlayerState = {
  track: null,
  queue: [],
  queueIndex: -1,
  paused: true,
  currentTime: 0,
  duration: 0,
  volume: parseFloat(localStorage.getItem("harmonic:volume") ?? "0.8"),
};
```

### Exported store and actions

```typescript
import { writable } from "svelte/store";

export const player = writable<PlayerState>(initial);
```

Each action is a plain exported function that calls `player.update(...)`:

#### `playTrack(track: AudioFile, queue: AudioFile[], index: number)`

Sets `track`, `queue`, `queueIndex`, resets `currentTime` and `duration` to 0, sets `paused` to `false`.

#### `togglePause()`

Flips `paused`. If there's no track loaded, no-op.

#### `nextTrack()`

If `queueIndex < queue.length - 1`, increment index and set the new track. Otherwise stop playback (set `track` to `null`, `paused` to `true`). Reset `currentTime` and `duration` to 0.

#### `prevTrack()`

If `currentTime > 3`, restart the current track (set `currentTime` to 0) — this is the standard music player UX where "previous" first restarts, then goes back. Otherwise, if `queueIndex > 0`, decrement and load the previous track.

#### `seek(time: number)`

Sets `currentTime`. The NowPlaying component will apply this to the `<audio>` element.

#### `setVolume(vol: number)`

Sets `volume`, clamped to `[0, 1]`. Writes to `localStorage` under key `harmonic:volume`.

### Important: the store doesn't own the `<audio>` element

The store is just data. The `NowPlaying` component (step 3) creates the `<audio>` element and reactively syncs it with the store. This separation keeps the store testable and avoids lifecycle issues.

### Helper: `audioUrl(track: AudioFile): string`

Export a pure helper that converts a track to a playable URL:

```typescript
export function audioUrl(track: AudioFile): string {
  return "/audio?path=" + encodeURIComponent(track.file_path);
}
```

---

## Step 3: Frontend — NowPlaying component

### New file: `frontend/src/lib/player/NowPlaying.svelte`

A fixed bar at the bottom of the viewport. Only renders when a track is loaded.

### Structure

```svelte
{#if $player.track}
<div class="now-playing" transition:fly={{ y: 64, duration: 200 }}>
  <audio ...></audio>

  <div class="np-info">      <!-- left: track info -->
  <div class="np-controls">  <!-- center: transport + seek -->
  <div class="np-volume">    <!-- right: volume -->
</div>
{/if}
```

### The `<audio>` element

This is the core. The element is hidden (no `controls` attribute). The component drives it reactively:

```svelte
<audio
  bind:this={audioEl}
  src={audioUrl($player.track)}
  bind:currentTime
  bind:duration
  bind:paused
  bind:volume
  on:ended={handleEnded}
  preload="metadata"
/>
```

Svelte's `bind:` makes this simple. But there are subtleties:

**Reactive sync between store and `<audio>`:**

The `<audio>` element is the source of truth for `currentTime` and `duration` (it fires `timeupdate` and `loadedmetadata` events). The component should write these values into the store so other components (like TrackRow's playing indicator) can read them.

When the store's `track` changes (user picked a new track), the `src` changes and the element auto-loads. You need to call `audioEl.play()` after the new source loads — use the `on:canplay` or `on:loadeddata` event, or reactively call play when `paused` is false and src changes.

When the user seeks via the seek bar, update `audioEl.currentTime` directly. The simplest pattern:

```svelte
<script>
  let audioEl: HTMLAudioElement;
  let localTime = 0;   // bound to <audio>.currentTime via timeupdate
  let localDur = 0;     // bound via loadedmetadata
  let seeking = false;  // true while user is dragging the seek bar

  // Write audio time to store periodically
  function onTimeUpdate() {
    if (!seeking) {
      localTime = audioEl.currentTime;
      player.update(s => ({ ...s, currentTime: localTime }));
    }
  }

  // When store says "play a new track", handle the src swap
  $: if ($player.track && audioEl) {
    const newSrc = audioUrl($player.track);
    if (audioEl.src !== newSrc) {
      audioEl.src = newSrc;
      audioEl.load();
    }
  }

  // Sync play/pause state
  $: if (audioEl) {
    if ($player.paused) audioEl.pause();
    else audioEl.play().catch(() => {});
  }

  // Sync volume
  $: if (audioEl) {
    audioEl.volume = $player.volume;
  }

  function handleEnded() {
    nextTrack();
  }
</script>
```

**Why not just use `bind:currentTime` directly?** You can, but seeking gets tricky — when the user drags the range input, you don't want `timeupdate` events fighting with the drag. The `seeking` flag pattern avoids that.

### Left section: Track info

```svelte
<div class="np-info">
  <span class="np-title">{$player.track.title}</span>
  <span class="np-ext">
    <span class="ext-badge ext-{$player.track.ext}">
      {$player.track.ext.toUpperCase()}
    </span>
  </span>
</div>
```

Reuse the `.ext-badge` and `.ext-*` color classes from TrackRow. Either move them to a shared stylesheet or duplicate (they're ~10 lines).

### Center section: Transport controls + seek bar

**Transport buttons:**

```
[prev] [play/pause] [next]
```

Each is a `<button>` with an inline SVG icon. Use the same sizing/style approach as existing buttons in the codebase.

- Prev: calls `prevTrack()` from the store
- Play/pause: calls `togglePause()` from the store. Icon switches between play triangle and pause bars based on `$player.paused`.
- Next: calls `nextTrack()` from the store

**Seek bar:**

```svelte
<span class="time">{formatTime(localTime)}</span>
<input
  type="range"
  min="0"
  max={localDur}
  step="0.1"
  value={localTime}
  on:pointerdown={() => seeking = true}
  on:pointerup={commitSeek}
  on:input={onSeekInput}
/>
<span class="time">{formatTime(localDur)}</span>
```

`formatTime` converts seconds to `m:ss` format:

```typescript
function formatTime(s: number): string {
  if (!s || !isFinite(s)) return "0:00";
  const m = Math.floor(s / 60);
  const sec = Math.floor(s % 60);
  return `${m}:${sec.toString().padStart(2, "0")}`;
}
```

Seek interaction:
- `on:pointerdown` → set `seeking = true` to stop `timeupdate` from overriding the slider
- `on:input` → update the visual position of the slider (`localTime = e.target.value`) but don't seek yet
- `on:pointerup` → set `audioEl.currentTime = localTime`, set `seeking = false`

### Right section: Volume

A speaker icon (SVG) + a small range slider (0 to 1, step 0.01). Clicking the icon toggles mute (store the pre-mute volume so you can restore it).

```svelte
<button on:click={toggleMute}>
  <!-- swap between speaker / speaker-muted SVG based on volume === 0 -->
</button>
<input
  type="range"
  min="0"
  max="1"
  step="0.01"
  value={$player.volume}
  on:input={(e) => setVolume(parseFloat(e.target.value))}
/>
```

### Styling

The bar should:
- Be `position: fixed; bottom: 0; left: 0; right: 0;` with `height: 64px`
- Use `background: var(--surface)` with `border-top: 1px solid var(--border)`
- Use the three-column grid layout like the Toolbar: `grid-template-columns: 1fr auto 1fr`
- Seek bar accent color: `var(--primary)` for the filled portion (use the CSS `background: linear-gradient(...)` trick on the range input, or `accent-color: var(--primary)` if you just want it simple)
- Font sizes: `--text-sm` for title, `--text-xs` for time display

---

## Step 4: Wire Library → Player

### Modify: `frontend/src/lib/library/TrackRow.svelte`

**Add double-click to play:**

The component currently dispatches a `select` event on click. Add a second event for playback:

```svelte
on:dblclick={() => dispatch('play')}
```

Add `export let playing: boolean = false;` prop so the parent can indicate which track is currently playing.

**Visual playing indicator:**

The component already shows a play icon when `selected`. Change this: show the play icon when `playing` is true (regardless of selection). When both `playing` and `selected`, show the play icon with the accent color. A subtle pulsing animation on the icon is a nice touch but optional.

```svelte
<span class="col-num track-num">
  {#if playing}
    <svg ...class="playing-icon">...</svg>
  {:else}
    {index + 1}
  {/if}
</span>
```

### Modify: `frontend/src/views/Library.svelte`

**Import the store:**

```typescript
import { player, playTrack } from "../lib/stores/player";
```

**Handle play events from TrackRow:**

```svelte
<TrackRow
  ...
  playing={$player.track?.file_path === track.file_path}
  on:select={() => selectTrack(i)}
  on:play={() => playTrack(track, filteredTracks, i)}
/>
```

**Adjust layout height:**

The `.library` container currently uses `height: calc(100vh - 52px)` (viewport minus toolbar). When the player is visible, subtract the player bar height too:

```css
.library {
  height: calc(100vh - 52px);  /* no player */
}

.library.has-player {
  height: calc(100vh - 52px - 64px);  /* toolbar + player */
}
```

Toggle the class based on `$player.track !== null`. You'll need to import the store in Library.svelte for this (you're already importing it for the playing indicator).

### Modify: `frontend/src/App.svelte`

Mount the NowPlaying component:

```svelte
<script lang="ts">
  import Toolbar from "./lib/components/Toolbar.svelte";
  import Library from "./views/Library.svelte";
  import NowPlaying from "./lib/player/NowPlaying.svelte";
</script>

<Toolbar title="Harmonic" />
<main>
  <Library />
</main>
<NowPlaying />
```

`NowPlaying` handles its own visibility (renders nothing when no track is loaded).

---

## Step 5: Keyboard shortcuts

### Add to: `NowPlaying.svelte`

Use `<svelte:window>` to capture key events globally:

```svelte
<svelte:window on:keydown={handleKeydown} />
```

```typescript
function handleKeydown(e: KeyboardEvent) {
  // Don't capture when user is typing in an input
  if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

  switch (e.key) {
    case " ":
      e.preventDefault();  // prevent page scroll
      togglePause();
      break;
    case "ArrowLeft":
      seek(Math.max(0, $player.currentTime - 5));
      break;
    case "ArrowRight":
      seek(Math.min($player.duration, $player.currentTime + 5));
      break;
    case "ArrowUp":
      e.preventDefault();
      setVolume(Math.min(1, $player.volume + 0.05));
      break;
    case "ArrowDown":
      e.preventDefault();
      setVolume(Math.max(0, $player.volume - 0.05));
      break;
  }
}
```

---

## Implementation order

The steps above are in dependency order. Specifically:

1. **Step 1** (Go handler) — no frontend dependencies, can be tested standalone via devtools console
2. **Step 2** (player store) — no UI dependencies, pure TypeScript
3. **Step 3** (NowPlaying component) — depends on steps 1 + 2
4. **Step 4** (wire Library) — depends on steps 2 + 3
5. **Step 5** (keyboard shortcuts) — depends on step 3, trivially added to NowPlaying

Each step is independently testable before moving on.

---

## File summary

| File | Action | What |
|------|--------|------|
| `internal/fs/handler.go` | Create | Audio file HTTP handler using `http.ServeFile` |
| `main.go` | Modify | Add `Handler: fs.NewAudioHandler()` to `AssetServer` options |
| `frontend/src/lib/stores/player.ts` | Create | Writable store + action functions for player state |
| `frontend/src/lib/player/NowPlaying.svelte` | Create | Bottom bar: hidden `<audio>` element, transport controls, seek, volume |
| `frontend/src/App.svelte` | Modify | Mount `<NowPlaying />` |
| `frontend/src/views/Library.svelte` | Modify | Import store, pass `playing` prop, adjust height when player visible |
| `frontend/src/lib/library/TrackRow.svelte` | Modify | Add `dblclick` → play event, `playing` prop, playing indicator |

No new dependencies. No changes to `package.json` or `go.mod`.
