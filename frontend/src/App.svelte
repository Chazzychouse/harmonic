<script lang="ts">
  import { onMount } from "svelte";
  import Toolbar from "./lib/components/Toolbar.svelte";
  import Library from "./views/Library.svelte";
  import NowPlaying from "./lib/player/NowPlaying.svelte";
  import { initMedia } from "./lib/stores/player";

  onMount(() => { initMedia(); });
</script>

<Toolbar title="Harmonic" />
<main>
  <Library />
</main>
<NowPlaying />
