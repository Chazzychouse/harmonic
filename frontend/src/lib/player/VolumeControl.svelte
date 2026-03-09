<script lang="ts">
  export let volume: number;

  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher<{
    toggleMute: void;
    volumeInput: number;
  }>();

  function onVolumeInput(e: Event) {
    dispatch("volumeInput", parseFloat((e.target as HTMLInputElement).value));
  }
</script>

<div class="np-volume">
  <button class="ctrl-btn" on:click={() => dispatch("toggleMute")} title={volume === 0 ? "Unmute" : "Mute"}>
    {#if volume === 0}
      <!-- Speaker icon with X (muted) -->
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" fill="currentColor" stroke="none"/> <!-- Speaker cone -->
        <line x1="23" y1="9" x2="17" y2="15"/> <!-- X slash top-right to bottom-left -->
        <line x1="17" y1="9" x2="23" y2="15"/> <!-- X slash bottom-right to top-left -->
      </svg>
    {:else}
      <!-- Speaker icon with sound waves -->
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" fill="currentColor" stroke="none"/> <!-- Speaker cone -->
        {#if volume > 0.5}
          <path d="M19.07 4.93a10 10 0 0 1 0 14.14"/> <!-- Outer sound wave arc -->
        {/if}
        <path d="M15.54 8.46a5 5 0 0 1 0 7.07"/> <!-- Inner sound wave arc -->
      </svg>
    {/if}
  </button>
  <input
    class="volume-bar"
    type="range"
    min="0"
    max="1"
    step="0.01"
    value={volume}
    on:input={onVolumeInput}
  />
</div>

<style>
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

  .np-volume {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    justify-content: flex-end;
  }

  .volume-bar {
    -webkit-appearance: none;
    appearance: none;
    height: 3px;
    border-radius: 2px;
    background: var(--border-strong);
    outline: none;
    cursor: pointer;
    accent-color: var(--primary);
    width: 80px;
  }

  .volume-bar::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--primary);
    cursor: pointer;
  }
</style>
