<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let disabled: boolean = false;
  export let icon: boolean = false;
</script>

<button
  class="btn {variant} {size}"
  class:icon
  {disabled}
  on:click
>
  <slot />
</button>

<style>
  .btn {
    font-family: var(--font-family);
    border: 1px solid transparent;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: var(--transition);
    font-weight: 700;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    letter-spacing: 0.01em;
  }

  .btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  /* Sizes */
  .sm { padding: var(--space-1) var(--space-3); font-size: var(--text-xs); }
  .md { padding: var(--space-2) var(--space-4); font-size: var(--text-sm); }
  .lg { padding: var(--space-3) var(--space-6); font-size: var(--text-md); }

  /* Icon-only */
  .icon.sm { padding: var(--space-1); }
  .icon.md { padding: var(--space-2); }
  .icon.lg { padding: var(--space-3); }

  /* Primary — full gradient */
  .primary {
    background: var(--gradient);
    color: white;
    border-color: transparent;
  }
  .primary:hover:not(:disabled) {
    background: var(--gradient-h);
    box-shadow: 0 0 18px var(--gradient-glow);
  }

  /* Secondary */
  .secondary {
    background: var(--surface-2);
    color: var(--text-primary);
    border-color: var(--border);
  }
  .secondary:hover:not(:disabled) {
    background: var(--surface-hover);
    border-color: var(--border-strong);
  }

  /* Danger */
  .danger {
    background: var(--danger);
    color: white;
  }
  .danger:hover:not(:disabled) {
    background: var(--danger-hover);
    box-shadow: 0 0 18px rgba(239, 68, 68, 0.35);
  }

  /* Focus */
  .btn:focus-visible {
    outline: none;
    box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--primary);
  }
</style>
