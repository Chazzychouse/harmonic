<script lang="ts">
  export let title: string = '';
</script>

<header class="toolbar">
  <div class="left">
    {#if title}
      <span class="title">{title}</span>
    {/if}
    <slot name="left" />
  </div>
  <div class="center">
    <slot name="center" />
  </div>
  <div class="right">
    <slot name="right" />
  </div>
</header>

<style>
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-3) var(--space-4);
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .left, .center, .right {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  .left { flex: 1; }
  .center { flex: 1; justify-content: center; }
  .right { flex: 1; justify-content: flex-end; }

  .title {
    font-size: var(--text-md);
    font-weight: 700;
    color: var(--text-primary);
  }
</style>
