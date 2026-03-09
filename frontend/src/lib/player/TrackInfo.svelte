<script lang="ts">
  export let title: string;
  export let ext: string;
</script>

<div class="np-info">
  <span class="np-title">{title}</span>
  <span class="np-ext">
    <span class="ext-badge ext-{ext}">
      {ext.toUpperCase()}
    </span>
  </span>
</div>

<style>
  .np-info {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    overflow: hidden;
    min-width: 0;
  }

  .np-title {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .np-ext { flex-shrink: 0; }

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
