<script lang="ts">
  import { SelectFile, RemoveWatcher } from "../wailsjs/go/main/App";
  import { EventsOn } from "../wailsjs/runtime";

  interface Logs {
    [key: string]: string[];
  }

  let logs: Logs = {};
  let activeLog = "";

  let paused = false;
  const togglePaused = () => (paused = !paused);

  const handleLoadLog = async () => {
    const result = await SelectFile();
    if (result.success) {
      const path = result.data.path;
      logs[path] = result.data.lines;

      activeLog = path;

      setTimeout(scrollToBottom, 500);

      EventsOn(path, line => {
        const log = logs[path];
        logs[path] = [...log, line];

        if (!paused && path === activeLog) {
          scrollToBottom();
        }
      });
    } else {
      // TODO toast
    }
  };

  const scrollToBottom = () => {
    const linesContainer = document.getElementById("lines");
    linesContainer.scrollTo(0, linesContainer.scrollHeight);
  };

  const setActiveLog = (key: string) => {
    activeLog = key;

    logs = { ...logs };
  };

  const removeLog = (path: string) => {
    RemoveWatcher(path);

    delete logs[path];

    logs = { ...logs };

    if (path === activeLog) {
      const entries = Object.entries(logs);
      activeLog = entries.length >= 1 ? entries[0][0] : "";
    }
  };

  const getTabClass = (tab: string) =>
    tab === activeLog ? "tab tab-lifted tab-active" : "tab tab-lifted";
</script>

<main class="mx-4 h-full">
  <div class="mt-2 flex justify-between">
    <div>
      <button type="button" class="btn-ghost btn" aria-label="Open log" on:click={handleLoadLog}>
        <svg
          aria-hidden="true"
          class="h-6 w-6"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M3.75 9.776c.112-.017.227-.026.344-.026h15.812c.117 0 .232.009.344.026m-16.5 0a2.25 2.25 0 00-1.883 2.542l.857 6a2.25 2.25 0 002.227 1.932H19.05a2.25 2.25 0 002.227-1.932l.857-6a2.25 2.25 0 00-1.883-2.542m-16.5 0V6A2.25 2.25 0 016 3.75h3.879a1.5 1.5 0 011.06.44l2.122 2.12a1.5 1.5 0 001.06.44H18A2.25 2.25 0 0120.25 9v.776"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </button>

      <button
        type="button"
        class="btn-ghost btn ml-4"
        aria-label={paused ? "Play" : "Pause"}
        on:click={togglePaused}
      >
        {#if paused}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="h-6 w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z"
            />
          </svg>
        {:else}
          <svg
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="h-6 w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15.75 5.25v13.5m-7.5-13.5v13.5"
            />
          </svg>
        {/if}
      </button>
    </div>
  </div>

  <div class="mt-2">
    {#each Object.entries(logs) as [key] (key)}
      <div class={getTabClass(key)}>
        <button
          type="button"
          class="btn-ghost btn-xs btn normal-case"
          on:click={() => setActiveLog(key)}
        >
          {key.split(/[\\/]/).pop()}
        </button>

        <button class="btn-ghost btn-xs btn" on:click={() => removeLog(key)}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="h-6 w-6"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    {/each}
  </div>

  <div class="h-5/6 overflow-auto border-b" id="lines">
    {#if activeLog}
      {#each logs[activeLog] as line}
        <p class="whitespace-nowrap font-mono">{line}</p>
      {/each}
    {/if}
  </div>
</main>
