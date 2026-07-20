<script lang="ts">
  import { onMount } from 'svelte';
  import type { InitResponse } from './types/InitResponse';
  import type { QueryResponse } from './types/QueryResponse';

  let isLoading = true;
  let rdfFilePath = '';
  let queryText = '';
  let errorText = '';
  let queryResult: QueryResponse;

  let executeQueryDialog: HTMLDialogElement;

  onMount(async () => {
    try {
      const res = await fetch('http://localhost:4711/api/init'); // TODO do not hardcode backend uri

      if (res.ok) {
        let responseObj: InitResponse = await res.json();
        rdfFilePath = responseObj.rdfFilePath;
      }
    } catch (ex) {
      if (ex instanceof Error) {
        errorText = ex.message;
      } else {
        errorText = 'An error occured';
      }
    } finally {
      isLoading = false;
    }
  });

  const executeQuery = async () => {
    executeQueryDialog.showModal();
    const params = new URLSearchParams(); // TODO use multipart formdata
    params.append('query', queryText);

    try {
      const queryRes = await fetch('http://localhost:4711/api/query', {
        method: 'POST',
        body: params,
      });

      if (queryRes.ok) {
        queryResult = await queryRes.json();
      } else {
        errorText = queryRes.status + ' ' + (await queryRes.text());
      }
    } catch (ex) {
      if (ex instanceof Error) {
        errorText = ex.message;
      } else {
        errorText = 'An error occured';
      }
    } finally {
      executeQueryDialog.close();
    }
  };

  const cancelQuery = () => {
    // TODO cancel request
    executeQueryDialog.close();
  };
</script>

<section id="center">
  {#if isLoading}
    <div>Initializing ... Please wait</div>
  {:else}
    {#if errorText}
      <div>Error: {errorText}</div>
    {/if}
    <div class="top-message">
      You are using sparql2go with file
      <code class="file-path" title={rdfFilePath}>{rdfFilePath}</code>
    </div>

    <textarea placeholder="Enter query" bind:value={queryText}></textarea>
    <button on:click={executeQuery}>Execute</button>

    {#if queryResult}
      Result: {queryResult.count} rows
    {/if}
  {/if}
</section>

<dialog bind:this={executeQueryDialog}>
  <h2>Executing query ...</h2>
  <button on:click={cancelQuery}>Cancel query ...</button>
</dialog>

<style>
  .top-message {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
  }

  .file-path {
    margin-left: 0.25rem;
    font-weight: 600;
  }
</style>
