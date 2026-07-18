<script lang="ts">
  import { onMount } from 'svelte';
  import type { InitResponse } from './types/InitResponse';

  let isLoading = true;
  let rdfFilePath = '';

  onMount(async () => {
    try {
      const res = await fetch('http://localhost:4711/api/init'); // TODO do not hardcode backend uri

      if (res.ok) {
        let responseObj: InitResponse = await res.json();
        rdfFilePath = responseObj.rdfFilePath;
      }
    } finally {
      isLoading = false;
    }
  });
</script>

<section id="center">
  {#if isLoading}
    <div>Initializing ... Please wait</div>
  {:else}
    <div class="top-message">You are using sparql2go with file "{rdfFilePath}"</div>
  {/if}
</section>

<style>
  .top-message {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
  }
</style>
