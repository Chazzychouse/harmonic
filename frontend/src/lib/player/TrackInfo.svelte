<script lang="ts">
  export let title: string;
  export let ext: string;
  export let artist: string = "";
  export let filePath: string = "";
  export let hasArt: boolean = false;
</script>

<div class="np-info">
  {#if hasArt}
    <img
      class="np-art"
      src="/art?path={encodeURIComponent(filePath)}"
      alt=""
    />
  {:else}
    <div class="np-art-placeholder">♪</div>
  {/if}
  <div class="np-text">
    <span class="np-title">{title}</span>
    <span class="np-artist">{artist || "Unknown Artist"}</span>
  </div>
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

  .np-art {
    width: 48px;
    height: 48px;
    border-radius: var(--radius-sm);
    object-fit: cover;
    flex-shrink: 0;
  }

  .np-art-placeholder {
    width: 48px;
    height: 48px;
    border-radius: var(--radius-sm);
    background: var(--surface-2);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    font-size: var(--text-lg);
    flex-shrink: 0;
  }

  .np-text {
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-width: 0;
    gap: 2px;
  }

  .np-title {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .np-artist {
    font-size: var(--text-xs);
    color: var(--text-secondary);
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
