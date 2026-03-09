<script lang="ts">
  import { formatTime } from "./formatTime";

  export let paused: boolean;
  export let currentTime: number;
  export let duration: number;

  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher<{
    prev: void;
    togglePause: void;
    next: void;
    seekInput: number;
    seekStart: void;
    seekCommit: void;
  }>();

  function onSeekInput(e: Event) {
    dispatch("seekInput", parseFloat((e.target as HTMLInputElement).value));
  }
</script>

<div class="np-controls">
  <div class="transport">
    <button class="ctrl-btn" on:click={() => dispatch("prev")} title="Previous">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
        <polygon points="19 20 9 12 19 4 19 20"/>
        <rect x="4" y="4" width="2" height="16" rx="1"/>
      </svg>
    </button>
    <button class="ctrl-btn play-btn" on:click={() => dispatch("togglePause")} title={paused ? "Play" : "Pause"}>
      {#if paused}
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
    <button class="ctrl-btn" on:click={() => dispatch("next")} title="Next">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
        <polygon points="5 4 15 12 5 20 5 4"/>
        <rect x="18" y="4" width="2" height="16" rx="1"/>
      </svg>
    </button>
  </div>
  <div class="seek-row">
    <span class="time">{formatTime(currentTime)}</span>
    <input
      class="seek-bar"
      type="range"
      min="0"
      max={duration || 1}
      step="0.1"
      value={currentTime}
      on:pointerdown={() => dispatch("seekStart")}
      on:pointerup={() => dispatch("seekCommit")}
      on:input={onSeekInput}
    />
    <span class="time right">{formatTime(duration)}</span>
  </div>
</div>

<style>
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

  .seek-bar {
    -webkit-appearance: none;
    appearance: none;
    height: 3px;
    border-radius: 2px;
    background: var(--border-strong);
    outline: none;
    cursor: pointer;
    accent-color: var(--primary);
    flex: 1;
  }

  .seek-bar::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--primary);
    cursor: pointer;
  }
</style>
