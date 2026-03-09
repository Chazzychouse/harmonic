<script lang="ts">
  import { player, prevTrack as storePrevTrack, togglePause, nextTrack, setVolume, toggleMute, audioUrl } from "../stores/player";
  import { fly } from "svelte/transition";
  import TrackInfo from "./TrackInfo.svelte";
  import TransportControls from "./TransportControls.svelte";
  import VolumeControl from "./VolumeControl.svelte";

  let audioEl: HTMLAudioElement;
  let localTime = 0;
  let localDur = 0;
  let seeking = false;

  function onTimeUpdate() {
    if (!seeking && audioEl) {
      localTime = audioEl.currentTime;
      player.update(s => ({ ...s, currentTime: localTime }));
    }
  }

  function onLoadedMetadata() {
    if (audioEl) {
      localDur = audioEl.duration;
      player.update(s => ({ ...s, duration: localDur }));
    }
  }

  function onAudioError() {
    const err = audioEl?.error;
    console.error("[audio] error code:", err?.code, "message:", err?.message, "src:", audioEl?.src);
  }

  $: paused = $player.paused;
  $: if (audioEl) {
    if (paused) audioEl.pause();
    else audioEl.play().catch(() => {});
  }
  $: if (audioEl) audioEl.volume = $player.volume;

  function handlePrev() {
    if (localTime > 3 && audioEl) {
      audioEl.currentTime = 0;
      localTime = 0;
      player.update(s => ({ ...s, currentTime: 0 }));
    } else {
      storePrevTrack();
    }
  }

  function handleSeekInput(e: CustomEvent<number>) {
    localTime = e.detail;
  }

  function handleSeekCommit() {
    if (audioEl) audioEl.currentTime = localTime;
    seeking = false;
    player.update(s => ({ ...s, currentTime: localTime }));
  }

  function seekTo(time: number) {
    if (audioEl) audioEl.currentTime = time;
    localTime = time;
    player.update(s => ({ ...s, currentTime: time }));
  }

  function handleKeydown(e: KeyboardEvent) {
    if (!$player.track) return;
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

    switch (e.key) {
      case " ":
        e.preventDefault();
        togglePause();
        break;
      case "ArrowLeft":
        e.preventDefault();
        seekTo(Math.max(0, localTime - 5));
        break;
      case "ArrowRight":
        e.preventDefault();
        seekTo(Math.min(localDur, localTime + 5));
        break;
      case "ArrowUp":
        e.preventDefault();
        setVolume(Math.min(1, $player.volume + 0.05));
        break;
      case "ArrowDown":
        e.preventDefault();
        setVolume(Math.max(0, $player.volume - 0.05));
        break;
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if $player.track}
<div class="now-playing" transition:fly={{ y: 64, duration: 200 }}>
  <audio
    bind:this={audioEl}
    src={audioUrl($player.track)}
    on:timeupdate={onTimeUpdate}
    on:loadedmetadata={onLoadedMetadata}
    on:ended={nextTrack}
    on:error={onAudioError}
    preload="metadata"
  ></audio>

  <TrackInfo
    title={$player.track.title}
    ext={$player.track.ext}
    artist={$player.track.artist}
    filePath={$player.track.file_path}
    hasArt={$player.track.has_art}
  />

  <TransportControls
    paused={$player.paused}
    currentTime={localTime}
    duration={localDur}
    on:prev={handlePrev}
    on:togglePause={togglePause}
    on:next={nextTrack}
    on:seekInput={handleSeekInput}
    on:seekStart={() => (seeking = true)}
    on:seekCommit={handleSeekCommit}
  />

  <VolumeControl
    volume={$player.volume}
    on:toggleMute={toggleMute}
    on:volumeInput={(e) => setVolume(e.detail)}
  />
</div>
{/if}

<style>
  .now-playing {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 64px;
    background: var(--surface);
    border-top: 1px solid var(--border);
    display: grid;
    grid-template-columns: 1fr auto 1fr;
    align-items: center;
    padding: 0 var(--space-4);
    gap: var(--space-4);
    z-index: 100;
  }
</style>
