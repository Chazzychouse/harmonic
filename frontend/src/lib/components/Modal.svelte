<script lang="ts">
  export let open: boolean = false;
  export let title: string = '';

  import { createEventDispatcher } from 'svelte';
  import { fade, scale } from 'svelte/transition';
  const dispatch = createEventDispatcher();

  function close() {
    dispatch('close');
  }

  function handleBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) close();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') close();
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if open}
  <div class="backdrop" on:click={handleBackdrop} on:keydown={handleKeydown} transition:fade={{ duration: 150 }}>
    <div class="modal" transition:scale={{ duration: 150, start: 0.97 }}>
      <header class="modal-header">
        <h2 class="modal-title">{title}</h2>
        <button class="close-btn" on:click={close}>×</button>
      </header>
      <div class="modal-body">
        <slot />
      </div>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.72);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 200;
  }

  .modal {
    background: var(--surface);
    border: 1px solid var(--border-strong);
    border-radius: var(--radius-lg);
    box-shadow: 0 32px 64px rgba(0, 0, 0, 0.6), 0 0 0 1px rgba(255,255,255,0.04);
    min-width: 420px;
    max-width: 90vw;
    max-height: 80vh;
    overflow: auto;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-4) var(--space-5);
    border-bottom: 1px solid var(--border);
  }

  .modal-title {
    margin: 0;
    font-size: var(--text-lg);
    font-weight: 700;
    color: var(--text-primary);
    letter-spacing: -0.01em;
  }

  .close-btn {
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: var(--text-xl);
    cursor: pointer;
    padding: 0;
    line-height: 1;
    transition: var(--transition);
  }

  .close-btn:hover {
    color: var(--text-primary);
  }

  .modal-body {
    padding: var(--space-5);
  }
</style>
