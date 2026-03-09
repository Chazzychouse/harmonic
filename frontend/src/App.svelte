<script lang="ts">
  import Toolbar from "./lib/components/Toolbar.svelte";
  import Library from "./views/Library.svelte";
  import NowPlaying from "./lib/player/NowPlaying.svelte";
</script>

<Toolbar title="Harmonic" />
<main>
  <Library />
</main>
<NowPlaying />
