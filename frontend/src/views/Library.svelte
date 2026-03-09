<script lang="ts">
  import { fade } from "svelte/transition";
  import FolderPickerModal from "../lib/library/FolderPickerModal.svelte";
  import EmptyState from "../lib/library/EmptyState.svelte";
  import ScanningState from "../lib/library/ScanningState.svelte";
  import ErrorState from "../lib/library/ErrorState.svelte";
  import TrackListToolbar from "../lib/library/TrackListToolbar.svelte";
  import TrackRow from "../lib/library/TrackRow.svelte";
  import { ScanDirectory } from "../../wailsjs/go/fs/FsService";
  import { player, playTrack } from "../lib/stores/player";

  interface AudioFile {
    title: string;
    file_path: string;
    ext: string;
  }

  let folderPath = "";
  let tracks: AudioFile[] = [];
  let filteredTracks: AudioFile[] = [];
  let selectedIndex: number | null = null;
  let searchQuery = "";
  let scanning = false;
  let scanError = "";
  let pickerOpen = false;

  $: {
    const q = searchQuery.toLowerCase();
    filteredTracks = q
      ? tracks.filter(
          (t) =>
            t.title.toLowerCase().includes(q) ||
            t.ext.toLowerCase().includes(q),
        )
      : tracks;
  }

  async function openLibrary(path: string) {
    pickerOpen = false;
    folderPath = path;
    scanning = true;
    scanError = "";
    tracks = [];
    selectedIndex = null;
    try {
      const results = await ScanDirectory(folderPath);
      tracks = results ?? [];
    } catch (e: unknown) {
      scanError = e instanceof Error ? e.message : String(e);
    } finally {
      scanning = false;
    }
  }

  function openPicker() {
    pickerOpen = true;
  }

  function closePicker() {
    pickerOpen = false;
  }

  function selectTrack(i: number) {
    selectedIndex = i === selectedIndex ? null : i;
  }

  function basename(path: string) {
    return path.split(/[\\/]/).pop() ?? path;
  }
</script>

<FolderPickerModal
  open={pickerOpen}
  initialPath={folderPath}
  on:select={(e) => openLibrary(e.detail)}
  on:close={closePicker}
/>

<!-- Main Library View -->
<div class="library" class:has-player={$player.track !== null}>
  <div class="library-main">
    {#if !folderPath && !scanning}
      <EmptyState on:add={openPicker} />
    {:else if scanning}
      <ScanningState folderName={basename(folderPath)} />
    {:else if scanError}
      <ErrorState message={scanError} on:retry={openPicker} />
    {:else}
      <div class="list-view" in:fade={{ duration: 200 }}>
        <TrackListToolbar
          {folderPath}
          trackCount={filteredTracks.length}
          {searchQuery}
          on:search={(e) => (searchQuery = e.detail)}
          on:changefolder={openPicker}
        />

        <!-- Column Headers -->
        <div class="col-header">
          <span class="col-num">#</span>
          <span class="col-title">Title</span>
          <span class="col-path">Path</span>
          <span class="col-ext">Format</span>
        </div>

        <!-- Tracks -->
        <div class="track-rows" role="listbox" aria-label="Track list">
          {#if filteredTracks.length === 0}
            <div class="no-results" in:fade={{ duration: 150 }}>
              No tracks match "<em>{searchQuery}</em>"
            </div>
          {:else}
            {#each filteredTracks as track, i (track.file_path)}
              <TrackRow
                index={i}
                title={track.title}
                filePath={track.file_path}
                ext={track.ext}
                selected={selectedIndex === i}
                playing={$player.track?.file_path === track.file_path}
                even={i % 2 === 0}
                delay={Math.min(i * 8, 300)}
                on:select={() => selectTrack(i)}
                on:play={() => playTrack(track, filteredTracks, i)}
              />
            {/each}
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .library {
    height: calc(100vh - 52px);
    transition: height 200ms ease;
    display: flex;
    flex-direction: row;
    overflow: hidden;
  }

  .library.has-player {
    height: calc(100vh - 52px - 64px);
  }

  .library-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-width: 0;
  }

  .list-view {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-height: 0;
  }

  /* ── Column Header ── */
  .col-header {
    display: grid;
    grid-template-columns: 44px 1fr 1fr 80px;
    align-items: center;
    gap: 0;
    padding: var(--space-2) 0;
    border-bottom: 1px solid var(--border);
    background: var(--surface-2);
    flex-shrink: 0;
  }

  .col-header span {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-muted);
    padding: 0 var(--space-3);
    user-select: none;
  }

  .col-num {
    text-align: right;
    padding-right: var(--space-4) !important;
  }

  /* ── Track Rows Container ── */
  .track-rows {
    flex: 1;
    overflow-y: auto;
    min-height: 0;
    scrollbar-width: thin;
    scrollbar-color: var(--border-strong) transparent;
  }

  .track-rows::-webkit-scrollbar {
    width: 4px;
  }

  .track-rows::-webkit-scrollbar-thumb {
    background: var(--border-strong);
    border-radius: 2px;
  }

  .no-results {
    padding: var(--space-8);
    text-align: center;
    font-size: var(--text-sm);
    color: var(--text-secondary);
  }
</style>
