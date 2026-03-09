<script lang="ts">
  import { fly } from "svelte/transition";
  import { createEventDispatcher } from "svelte";

  export let index: number;
  export let title: string;
  export let artist: string = "";
  export let filePath: string;
  export let ext: string;
  export let selected: boolean = false;
  export let playing: boolean = false;
  export let even: boolean = false;
  export let delay: number = 0;

  const dispatch = createEventDispatcher();

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" || e.key === " ") dispatch('select');
  }
</script>

<div
  class="track-row"
  class:selected
  class:even
  role="option"
  aria-selected={selected}
  tabindex="0"
  on:click={() => dispatch('select')}
  on:dblclick={() => dispatch('play')}
  on:keydown={handleKeydown}
  in:fly={{ y: 4, duration: 120, delay }}
>
  <span class="col-num track-num">
    {#if playing}
      <svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor" class="playing-icon" class:accent={selected}>
        <polygon points="5 3 19 12 5 21 5 3"/>
      </svg>
    {:else}
      {index + 1}
    {/if}
  </span>
  <span class="col-title track-title">{title}</span>
  <span class="col-artist track-artist">{artist || "—"}</span>
  <span class="col-ext track-ext">
    <span class="ext-badge ext-{ext}">{ext.toUpperCase()}</span>
  </span>
</div>

<style>
  .track-row {
    display: grid;
    grid-template-columns: 44px 1fr 1fr 80px;
    align-items: center;
    gap: 0;
    padding: 0;
    height: 38px;
    cursor: pointer;
    transition: background 100ms ease;
    border-bottom: 1px solid transparent;
    border-left: 2px solid transparent;
    outline: none;
  }

  .track-row.even {
    background: rgba(255, 255, 255, 0.015);
  }

  .track-row:hover {
    background: color-mix(in srgb, var(--primary) 10%, transparent);
  }

  .track-row:focus-visible {
    box-shadow: inset 0 0 0 1px var(--primary);
  }

  .track-row.selected {
    background: color-mix(in srgb, var(--primary) 18%, transparent) !important;
    border-left-color: var(--accent);
    border-bottom-color: color-mix(in srgb, var(--primary) 25%, transparent);
  }

  .track-row.selected .track-title {
    color: var(--text-primary);
  }

  .track-row.selected .track-num {
    color: var(--accent);
  }

  .track-row span {
    padding: 0 var(--space-3);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .track-num {
    text-align: right;
    font-size: var(--text-xs);
    color: var(--text-muted);
    padding-right: var(--space-4) !important;
    display: flex !important;
    align-items: center;
    justify-content: flex-end;
  }

  .playing-icon {
    color: var(--text-secondary);
  }

  .playing-icon.accent {
    color: var(--accent);
  }

  .track-title {
    font-size: var(--text-sm);
    color: var(--text-primary);
    font-weight: 500;
  }

  .track-artist {
    font-size: var(--text-xs);
    color: var(--text-secondary);
  }

  .col-ext {
    display: flex !important;
    align-items: center;
    justify-content: flex-start;
  }

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
</style>
