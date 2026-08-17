<script lang="ts">
  import { onMount } from 'svelte';
  import type { InitResponse } from './types/InitResponse';
  import type { QueryResponse } from './types/QueryResponse';

  let isLoading = $state(true);
  let rdfFilePath = $state('');
  let queryText = $state('');
  let errorText = $state('');
  let queryResult: QueryResponse | undefined = $state();
  let canExecuteQuery = $derived(!!queryText && queryText.length > 0);

  let executeQueryDialog: HTMLDialogElement;

  onMount(async () => {
    try {
      const res = await fetch('/api/init');

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

    try {
      const queryRes = await fetch('/api/query', {
        method: 'POST',
        body: queryText,
        headers: {
          'Content-Type': 'text/plain',
        },
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

<section id="center" class="query-workspace">
  {#if isLoading}
    <div class="status-message">Initializing ... Please wait</div>
  {:else}
    {#if errorText}
      <div class="error-message" role="alert">Error: {errorText}</div>
    {/if}
    <div class="top-message">
      You are using sparql2go with file
      <code class="file-path" title={rdfFilePath}>{rdfFilePath}</code>
    </div>

    <div class="query-panel">
      <label for="query">SPARQL query</label>
      <textarea id="query" placeholder="Enter query" bind:value={queryText}></textarea>
      <div class="actions">
        <button class="execute-button" onclick={executeQuery} disabled={!canExecuteQuery}
          >Execute query</button
        >
      </div>
    </div>

    {#if queryResult}
      <section class="results" aria-label="Query results">
        <div class="results-heading">
          <h2>Results</h2>
          <span>{queryResult.count} {queryResult.count === 1 ? 'row' : 'rows'}</span>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                {#each queryResult.vars as col}
                  <th>{col}</th>
                {/each}
              </tr>
            </thead>
            <tbody>
              {#each queryResult.result as row}
                <tr>
                  {#each queryResult.vars as col}
                    <td>
                      {row[col]}
                    </td>
                  {/each}
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </section>
    {/if}
  {/if}
</section>

<dialog bind:this={executeQueryDialog}>
  <h2>Executing query ...</h2>
  <p>Searching the RDF data. This may take a moment.</p>
  <button class="cancel-button" onclick={cancelQuery}>Cancel query</button>
</dialog>

<style>
  .query-workspace {
    width: min(100% - 48px, 900px);
    margin: 0 auto;
    padding: 52px 0;
    align-items: stretch;
    gap: 22px;
  }

  .status-message,
  .top-message,
  .error-message,
  .query-panel,
  .results {
    box-sizing: border-box;
    width: 100%;
  }

  .status-message,
  .top-message {
    color: var(--text);
    font-size: 0.9rem;
    text-align: left;
  }

  .error-message {
    padding: 12px 14px;
    border: 1px solid #e9a6a6;
    border-radius: 8px;
    color: #9b2424;
    background: #fff4f4;
    text-align: left;
  }

  .file-path {
    margin-left: 0.25rem;
    font-weight: 600;
    max-width: min(100%, 560px);
    overflow: hidden;
    vertical-align: middle;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .query-panel,
  .results {
    border: 1px solid var(--border);
    border-radius: 12px;
    background: color-mix(in srgb, var(--bg) 96%, var(--code-bg));
    box-shadow: 0 8px 24px rgba(23, 20, 29, 0.05);
  }

  .query-panel {
    padding: 20px;
  }

  label {
    display: block;
    margin-bottom: 10px;
    color: var(--text-h);
    font-size: 0.9rem;
    font-weight: 600;
    text-align: left;
  }

  textarea {
    box-sizing: border-box;
    display: block;
    width: 100%;
    min-height: 220px;
    resize: vertical;
    padding: 14px;
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--text-h);
    background: var(--bg);
    font: 0.9rem/1.55 var(--mono);
    outline: none;
  }

  textarea::placeholder {
    color: var(--text);
  }

  textarea:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px var(--accent-bg);
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 14px;
  }

  .execute-button,
  .cancel-button {
    padding: 9px 15px;
    border: 0;
    border-radius: 7px;
    cursor: pointer;
    font: 600 0.9rem var(--sans);
  }

  .execute-button {
    color: white;
    background: var(--accent);
  }

  .execute-button:not(:disabled):hover {
    filter: brightness(0.94);
  }

  .execute-button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .execute-button:focus-visible,
  .cancel-button:focus-visible {
    outline: 3px solid var(--accent-border);
    outline-offset: 2px;
  }

  .results {
    overflow: hidden;
  }

  .results-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border);
  }

  .results-heading h2 {
    margin: 0;
    font-size: 1rem;
  }

  .results-heading span {
    color: var(--text);
    font-size: 0.85rem;
  }

  .table-wrap {
    /* Keep large result sets self-contained so both scrollbars stay visible. */
    max-height: min(60vh, 640px);
    overflow: auto;
    overscroll-behavior: contain;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
  }

  th,
  td {
    padding: 12px 20px;
    border-bottom: 1px solid var(--border);
    white-space: nowrap;
  }

  th {
    position: sticky;
    top: 0;
    z-index: 1;
    color: var(--text-h);
    background: color-mix(in srgb, var(--code-bg) 70%, transparent);
    font-size: 0.8rem;
    font-weight: 650;
  }

  td {
    color: var(--text);
    font: 0.85rem/1.4 var(--mono);
  }

  tbody tr:last-child td {
    border-bottom: 0;
  }

  tbody tr:hover {
    background: var(--accent-bg);
  }

  dialog {
    width: min(360px, calc(100% - 48px));
    padding: 24px;
    border: 1px solid var(--border);
    border-radius: 12px;
    color: var(--text);
    background: var(--bg);
    box-shadow: var(--shadow);
  }

  dialog::backdrop {
    background: rgba(8, 6, 13, 0.3);
  }

  dialog p {
    margin: 0 0 20px;
    font-size: 0.9rem;
  }

  .cancel-button {
    color: var(--text-h);
    border: 1px solid var(--border);
    background: transparent;
  }

  .cancel-button:hover {
    background: var(--code-bg);
  }

  @media (max-width: 640px) {
    .query-workspace {
      width: min(100% - 32px, 900px);
      padding: 32px 0;
    }

    .query-panel {
      padding: 14px;
    }

    .results-heading,
    th,
    td {
      padding-inline: 14px;
    }
  }
</style>
