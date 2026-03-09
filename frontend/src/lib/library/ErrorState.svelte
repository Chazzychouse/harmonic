<script lang="ts">
  import { fade } from "svelte/transition";
  import { createEventDispatcher } from "svelte";
  import Button from "../components/Button.svelte";

  export let message: string;

  const dispatch = createEventDispatcher();
</script>

<div class="error-state" in:fade={{ duration: 200 }}>
  <p class="error-msg">Failed to scan directory: {message}</p>
  <Button variant="primary" on:click={() => dispatch('retry')}>Try Again</Button>
</div>

<style>
  .error-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--space-4);
  }

  .error-msg {
    color: var(--danger);
    font-size: var(--text-sm);
  }

</style>
