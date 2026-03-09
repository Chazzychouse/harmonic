<script lang="ts">
  import { fade } from "svelte/transition";
  import { createEventDispatcher } from "svelte";
  import Button from "../components/Button.svelte";

  const dispatch = createEventDispatcher();
</script>

<div class="empty-state" in:fade={{ duration: 300 }}>
  <div class="empty-icon">
    <svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
      <path d="M9 18V5l12-2v13"/>
      <circle cx="6" cy="18" r="3"/>
      <circle cx="18" cy="16" r="3"/>
    </svg>
  </div>
  <h2 class="empty-title">Your library is empty</h2>
  <p class="empty-subtitle">
    Choose a folder to discover your audio files
  </p>
  <Button variant="primary" size="lg" on:click={() => dispatch('add')}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
      <line x1="12" y1="11" x2="12" y2="17"/>
      <line x1="9" y1="14" x2="15" y2="14"/>
    </svg>
    Add Folder
  </Button>
</div>

<style>
  .empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--space-4);
    padding: var(--space-8);
  }

  .empty-icon {
    color: rgba(255, 255, 255, 0.12);
  }

  .empty-title {
    margin: 0;
    font-size: var(--text-xl);
    font-weight: 700;
    color: var(--text-primary);
    letter-spacing: -0.02em;
  }

  .empty-subtitle {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--text-secondary);
  }

</style>
