<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import Modal from "../components/Modal.svelte";
  import { SelectDirectory } from "../../../wailsjs/go/fs/FsService";

  export let open: boolean = false;
  export let initialPath: string = "";

  let pendingPath = initialPath;
  $: if (open) pendingPath = initialPath;

  const dispatch = createEventDispatcher<{ select: string; close: void }>();

  async function browseForFolder() {
    const dir = await SelectDirectory();
    if (dir) pendingPath = dir;
  }

  function confirm() {
    if (!pendingPath) return;
    dispatch("select", pendingPath);
  }

  function close() {
    dispatch("close");
  }
</script>

<Modal {open} title="Open Library Folder" on:close={close}>
  <div class="picker-body">
    <div class="picker-icon">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M3 7a2 2 0 0 1 2-2h3.93a2 2 0 0 1 1.66.9l.82 1.2A2 2 0 0 0 13.07 8H19a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z"/>
      </svg>
    </div>

    <p class="picker-description">
      Select a folder containing your music files. Harmonic will scan for
      all audio tracks inside, including subdirectories.
    </p>

    <div class="path-row">
      <div class="path-display" class:has-path={!!pendingPath}>
        {#if pendingPath}
          <span class="path-text">{pendingPath}</span>
        {:else}
          <span class="path-placeholder">No folder selected</span>
        {/if}
      </div>
      <button class="browse-btn" on:click={browseForFolder}>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="17 8 12 3 7 8"/>
          <line x1="12" y1="3" x2="12" y2="15"/>
        </svg>
        Browse
      </button>
    </div>

    <div class="picker-actions">
      <button class="action-btn cancel" on:click={close}>Cancel</button>
      <button
        class="action-btn confirm"
        disabled={!pendingPath}
        on:click={confirm}
      >
        Open Library
      </button>
    </div>
  </div>
</Modal>

<style>
  .picker-body {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
    align-items: center;
    text-align: center;
  }

  .picker-icon {
    color: var(--accent);
    opacity: 0.8;
  }

  .picker-description {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--text-secondary);
    line-height: 1.6;
    max-width: 360px;
  }

  .path-row {
    display: flex;
    gap: var(--space-2);
    width: 100%;
    align-items: stretch;
  }

  .path-display {
    flex: 1;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    padding: var(--space-2) var(--space-3);
    font-size: var(--text-xs);
    min-height: 36px;
    display: flex;
    align-items: center;
    text-align: left;
    overflow: hidden;
  }

  .path-display.has-path {
    border-color: var(--primary);
  }

  .path-text {
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .path-placeholder {
    color: var(--text-secondary);
    font-style: italic;
  }

  .browse-btn {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    background: var(--surface-2);
    border: 1px solid var(--border-strong);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-family: var(--font-family);
    font-size: var(--text-xs);
    font-weight: 600;
    padding: var(--space-2) var(--space-3);
    cursor: pointer;
    white-space: nowrap;
    transition: var(--transition);
    flex-shrink: 0;
  }

  .browse-btn:hover {
    background: var(--surface-hover);
    border-color: var(--primary);
  }

  .picker-actions {
    display: flex;
    gap: var(--space-3);
    justify-content: flex-end;
    width: 100%;
    padding-top: var(--space-2);
    border-top: 1px solid var(--border);
    margin-top: var(--space-2);
  }

  .action-btn {
    padding: var(--space-2) var(--space-5);
    border-radius: var(--radius-sm);
    font-family: var(--font-family);
    font-size: var(--text-sm);
    font-weight: 600;
    cursor: pointer;
    transition: var(--transition);
    border: 1px solid transparent;
  }

  .action-btn.cancel {
    background: var(--surface-2);
    border-color: var(--border);
    color: var(--text-secondary);
  }

  .action-btn.cancel:hover {
    background: var(--surface-hover);
    color: var(--text-primary);
    border-color: var(--border-strong);
  }

  .action-btn.confirm {
    background: var(--gradient);
    color: #fff;
    box-shadow: 0 4px 16px var(--gradient-glow);
  }

  .action-btn.confirm:hover:not(:disabled) {
    background: var(--gradient-h);
    box-shadow: 0 6px 24px var(--gradient-glow);
  }

  .action-btn.confirm:disabled {
    opacity: 0.35;
    cursor: not-allowed;
    box-shadow: none;
  }
</style>
