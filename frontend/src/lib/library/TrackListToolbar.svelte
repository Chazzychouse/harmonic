<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let folderPath: string;
  export let trackCount: number;
  export let searchQuery: string;

  const dispatch = createEventDispatcher();
</script>

<div class="list-toolbar">
  <div class="folder-info">
    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="folder-icon-sm">
      <path d="M3 7a2 2 0 0 1 2-2h3.93a2 2 0 0 1 1.66.9l.82 1.2A2 2 0 0 0 13.07 8H19a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z"/>
    </svg>
    <span class="folder-path">{folderPath}</span>
    <span class="track-count">{trackCount} track{trackCount !== 1 ? "s" : ""}</span>
  </div>
  <div class="toolbar-right">
    <div class="search-wrap">
      <svg class="search-icon" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"/>
        <line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
      <input
        class="search-input"
        type="text"
        placeholder="Search tracks…"
        value={searchQuery}
        on:input={(e) => dispatch('search', e.currentTarget.value)}
      />
    </div>
    <button class="change-folder-btn" on:click={() => dispatch('changefolder')}>
      Change Folder
    </button>
  </div>
</div>

<style>
  .list-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-3) var(--space-5);
    border-bottom: 1px solid var(--border);
    background: var(--surface);
    flex-shrink: 0;
    gap: var(--space-4);
  }

  .folder-info {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    min-width: 0;
    flex: 1;
  }

  .folder-icon-sm {
    color: var(--primary);
    flex-shrink: 0;
  }

  .folder-path {
    font-size: var(--text-xs);
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
  }

  .track-count {
    font-size: var(--text-xs);
    color: var(--text-muted);
    white-space: nowrap;
    flex-shrink: 0;
  }

  .toolbar-right {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    flex-shrink: 0;
  }

  .search-wrap {
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-icon {
    position: absolute;
    left: var(--space-3);
    color: var(--text-secondary);
    pointer-events: none;
  }

  .search-input {
    background: var(--surface-2);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-family: var(--font-family);
    font-size: var(--text-xs);
    padding: var(--space-2) var(--space-3) var(--space-2) 30px;
    width: 200px;
    outline: none;
    transition: var(--transition);
  }

  .search-input::placeholder {
    color: var(--text-secondary);
  }

  .search-input:focus {
    border-color: var(--primary);
    box-shadow: 0 0 0 2px var(--primary-glow);
  }

  .change-folder-btn {
    background: var(--surface-2);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    font-family: var(--font-family);
    font-size: var(--text-xs);
    padding: var(--space-2) var(--space-3);
    cursor: pointer;
    transition: var(--transition);
    white-space: nowrap;
  }

  .change-folder-btn:hover {
    color: var(--text-primary);
    background: var(--surface-hover);
    border-color: var(--border-strong);
  }
</style>
