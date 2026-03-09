<script lang="ts">
  import { player, prevTrack as storePrevTrack, togglePause, nextTrack, setVolume, toggleMute, audioUrl } from "../stores/player";
  import { fly } from "svelte/transition";

  let audioEl: HTMLAudioElement;
  let localTime = 0;
  let localDur = 0;
  let seeking = false;

  function onTimeUpdate() {
    if (!seeking && audioEl) {
      localTime = audioEl.currentTime;
      player.update(s => ({ ...s, currentTime: localTime }));
    }
  }

  function onLoadedMetadata() {
    if (audioEl) {
      localDur = audioEl.duration;
      player.update(s => ({ ...s, duration: localDur }));
    }
  }

  // Sync play/pause from store to audio element
  $: paused = $player.paused;
  $: if (audioEl) {
    if (paused) audioEl.pause();
    else audioEl.play().catch(() => {});
  }

  // Sync volume from store to audio element
  $: if (audioEl) audioEl.volume = $player.volume;

  function handleEnded() {
    nextTrack();
  }

  function onSeekInput(e: Event) {
    localTime = parseFloat((e.target as HTMLInputElement).value);
  }

  function commitSeek() {
    if (audioEl) audioEl.currentTime = localTime;
    seeking = false;
    player.update(s => ({ ...s, currentTime: localTime }));
  }

  function handlePrev() {
    if (localTime > 3 && audioEl) {
      audioEl.currentTime = 0;
      localTime = 0;
      player.update(s => ({ ...s, currentTime: 0 }));
    } else {
      storePrevTrack();
    }
  }

  function formatTime(s: number): string {
    if (!s || !isFinite(s)) return "0:00";
    const m = Math.floor(s / 60);
    const sec = Math.floor(s % 60);
    return `${m}:${sec.toString().padStart(2, "0")}`;
  }

  function onVolumeInput(e: Event) {
    setVolume(parseFloat((e.target as HTMLInputElement).value));
  }

  function handleKeydown(e: KeyboardEvent) {
    if (!$player.track) return;
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

    switch (e.key) {
      case " ":
        e.preventDefault();
        togglePause();
        break;
      case "ArrowLeft": {
        e.preventDefault();
        const t = Math.max(0, localTime - 5);
        if (audioEl) audioEl.currentTime = t;
        localTime = t;
        player.update(s => ({ ...s, currentTime: t }));
        break;
      }
      case "ArrowRight": {
        e.preventDefault();
        const t = Math.min(localDur, localTime + 5);
        if (audioEl) audioEl.currentTime = t;
        localTime = t;
        player.update(s => ({ ...s, currentTime: t }));
        break;
      }
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
</script>

<svelte:window on:keydown={handleKeydown} />

{#if $player.track}
<div class="now-playing" transition:fly={{ y: 64, duration: 200 }}>
  <audio
    bind:this={audioEl}
    src={audioUrl($player.track)}
    on:timeupdate={onTimeUpdate}
    on:loadedmetadata={onLoadedMetadata}
    on:ended={handleEnded}
    preload="metadata"
  ></audio>

  <!-- Left: track info -->
  <div class="np-info">
    <span class="np-title">{$player.track.title}</span>
    <span class="np-ext">
      <span class="ext-badge ext-{$player.track.ext}">
        {$player.track.ext.toUpperCase()}
      </span>
    </span>
  </div>

  <!-- Center: transport + seek -->
  <div class="np-controls">
    <div class="transport">
      <button class="ctrl-btn" on:click={handlePrev} title="Previous">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="19 20 9 12 19 4 19 20"/>
          <rect x="4" y="4" width="2" height="16" rx="1"/>
        </svg>
      </button>
      <button class="ctrl-btn play-btn" on:click={togglePause} title={$player.paused ? "Play" : "Pause"}>
        {#if $player.paused}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
        {:else}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <rect x="6" y="4" width="4" height="16" rx="1"/>
            <rect x="14" y="4" width="4" height="16" rx="1"/>
          </svg>
        {/if}
      </button>
      <button class="ctrl-btn" on:click={nextTrack} title="Next">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="5 4 15 12 5 20 5 4"/>
          <rect x="18" y="4" width="2" height="16" rx="1"/>
        </svg>
      </button>
    </div>
    <div class="seek-row">
      <span class="time">{formatTime(localTime)}</span>
      <input
        class="seek-bar"
        type="range"
        min="0"
        max={localDur || 1}
        step="0.1"
        value={localTime}
        on:pointerdown={() => (seeking = true)}
        on:pointerup={commitSeek}
        on:input={onSeekInput}
      />
      <span class="time right">{formatTime(localDur)}</span>
    </div>
  </div>

  <!-- Right: volume -->
  <div class="np-volume">
    <button class="ctrl-btn" on:click={toggleMute} title={$player.volume === 0 ? "Unmute" : "Mute"}>
      {#if $player.volume === 0}
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" fill="currentColor" stroke="none"/>
          <line x1="23" y1="9" x2="17" y2="15"/>
          <line x1="17" y1="9" x2="23" y2="15"/>
        </svg>
      {:else}
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" fill="currentColor" stroke="none"/>
          {#if $player.volume > 0.5}
            <path d="M19.07 4.93a10 10 0 0 1 0 14.14"/>
          {/if}
          <path d="M15.54 8.46a5 5 0 0 1 0 7.07"/>
        </svg>
      {/if}
    </button>
    <input
      class="volume-bar"
      type="range"
      min="0"
      max="1"
      step="0.01"
      value={$player.volume}
      on:input={onVolumeInput}
    />
  </div>
</div>
{/if}

<style>
  .now-playing {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 64px;
    background: var(--surface);
    border-top: 1px solid var(--border);
    display: grid;
    grid-template-columns: 1fr auto 1fr;
    align-items: center;
    padding: 0 var(--space-4);
    gap: var(--space-4);
    z-index: 100;
  }

  /* Left */
  .np-info {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    overflow: hidden;
    min-width: 0;
  }

  .np-title {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .np-ext { flex-shrink: 0; }

  .ext-badge {
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.06em;
    padding: 2px 6px;
    border-radius: 3px;
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.08);
  }

  .ext-mp3  { color: #60a5fa; border-color: rgba(96, 165, 250, 0.3); background: rgba(96, 165, 250, 0.08); }
  .ext-flac { color: #34d399; border-color: rgba(52, 211, 153, 0.3); background: rgba(52, 211, 153, 0.08); }
  .ext-wav  { color: #f59e0b; border-color: rgba(245, 158, 11, 0.3); background: rgba(245, 158, 11, 0.08); }
  .ext-aiff,
  .ext-aif  { color: #f59e0b; border-color: rgba(245, 158, 11, 0.3); background: rgba(245, 158, 11, 0.08); }
  .ext-m4a  { color: #f472b6; border-color: rgba(244, 114, 182, 0.3); background: rgba(244, 114, 182, 0.08); }
  .ext-ogg,
  .ext-opus { color: #a78bfa; border-color: rgba(167, 139, 250, 0.3); background: rgba(167, 139, 250, 0.08); }

  /* Center */
  .np-controls {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    min-width: 280px;
  }

  .transport {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  .ctrl-btn {
    background: none;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 4px;
    border-radius: 4px;
    transition: color 80ms, background 80ms;
  }

  .ctrl-btn:hover {
    color: var(--text-primary);
    background: rgba(255, 255, 255, 0.06);
  }

  .play-btn {
    color: var(--text-primary);
    background: rgba(255, 255, 255, 0.08);
    border-radius: 50%;
    padding: 7px;
  }

  .play-btn:hover { background: rgba(255, 255, 255, 0.14); }

  .seek-row {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    width: 100%;
  }

  .time {
    font-size: var(--text-xs);
    color: var(--text-muted);
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
    min-width: 32px;
  }

  .time.right { text-align: right; }

  .seek-bar,
  .volume-bar {
    -webkit-appearance: none;
    appearance: none;
    height: 3px;
    border-radius: 2px;
    background: var(--border-strong);
    outline: none;
    cursor: pointer;
    accent-color: var(--primary);
  }

  .seek-bar { flex: 1; }

  .seek-bar::-webkit-slider-thumb,
  .volume-bar::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--primary);
    cursor: pointer;
  }

  /* Right */
  .np-volume {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    justify-content: flex-end;
  }

  .volume-bar { width: 80px; }
</style>
