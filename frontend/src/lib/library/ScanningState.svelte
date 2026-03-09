<script lang="ts">
  import { fade } from "svelte/transition";

  export let folderName: string;
</script>

<div class="scanning-state" in:fade={{ duration: 200 }}>
  <div class="scan-spinner"></div>
  <p class="scan-text">Scanning <span class="scan-path">{folderName}</span>…</p>
</div>

<style>
  .scanning-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--space-4);
  }

  .scan-spinner {
    width: 36px;
    height: 36px;
    border: 2px solid color-mix(in srgb, var(--primary) 20%, transparent);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .scan-text {
    font-size: var(--text-sm);
    color: var(--text-secondary);
    margin: 0;
  }

  .scan-path {
    color: var(--accent);
  }
</style>
