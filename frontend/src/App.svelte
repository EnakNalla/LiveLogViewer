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

  let searchTerm = "";
  let searchResults = 0;
  let searchIndex = 0;

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
    tab === activeLog ? "tab tab-bordered tab-lifted tab-active" : "tab tab-lifted";

  const handleSearch = () => {
    // TODO handle clear
    if (!searchTerm) {
      clearSearch();
      return;
    }

    paused = true;
    let results = 0;
    const lines = document.querySelectorAll("p");

    for (let i = 0; i < lines.length; i++) {
      if (lines[i].innerText.includes(searchTerm)) {
        lines[i].innerHTML = lines[i].innerText.replace(
          searchTerm,
          '<span class="text-black bg-yellow-400" data-index=' +
            results +
            ">" +
            searchTerm +
            "</span>"
        );
        results++;
      }
    }

    searchResults = results - 1;

    if (searchResults <= 0) {
      // TODO toast
      clearSearch();
    }
  };

  const handleSearchInput = (e: Event) => {
    const target = e.target as HTMLInputElement;

    searchTerm = target.value;

    if (!searchTerm) {
      clearSearch();
    }
  };

  const clearSearch = () => {
    const highlights = document.querySelectorAll(".bg-yellow-400");
    highlights.forEach(element => element.classList.remove("bg-yellow-400", "text-black"));
    const prevResult = document.getElementsByClassName("current")[0];
    prevResult?.classList.remove("current", "bg-primary", "text-white");

    searchTerm = "";
    searchResults = 0;
    searchIndex = 0;
    paused = false;
  };

  const handleChevronClick = (next: boolean) => {
    searchIndex = next ? searchIndex + 1 : searchIndex - 1;

    if (searchIndex < 0) {
      searchIndex = searchResults;
    } else if (searchIndex >= searchResults) {
      searchIndex = 0;
    }

    focusResult(searchIndex);
  };

  const focusResult = (index: number) => {
    const prevResult = document.getElementsByClassName("current")[0];
    prevResult?.classList.remove("current", "bg-primary", "text-white");

    const result = document.querySelector('[data-index="' + index + '"]');
    result?.parentElement?.classList.add("bg-primary", "text-white", "current");
    result?.scrollIntoView();
  };
</script>

<main class="mx-4 h-full">
  <div class="flex justify-between pt-2">
    <div>
      <button type="button" class="btn-primary btn" aria-label="Open log" on:click={handleLoadLog}>
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
        class="btn-primary btn ml-8"
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

    <div class="flex">
      <label class="sr-only" for="search-input">Search</label>
      <form class="input-group" on:submit|preventDefault={handleSearch}>
        <div class="relative">
          <input
            value={searchTerm}
            on:change={handleSearchInput}
            type="text"
            class="input-bordered input rounded-r-none pr-8"
            placeholder="Enter search term"
            id="search-input"
          />
          <button
            class="btn-ghost btn-sm btn absolute bottom-2 right-2 p-0"
            type="button"
            on:click={() => clearSearch()}
          >
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
        <button type="submit" class="btn-primary btn">
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
              d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>
      </form>
      {#if searchResults}
        <p>{searchIndex} / {searchResults}</p>
        <div class="btn-group">
          <button
            type="button"
            class="btn btn-primary"
            aria-label="Next result"
            on:click={() => handleChevronClick(true)}
          >
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
                d="M19.5 5.25l-7.5 7.5-7.5-7.5m15 6l-7.5 7.5-7.5-7.5"
              />
            </svg>
          </button>
          <button
            type="button"
            class="btn btn-primary"
            aria-label="Previous result"
            on:click={() => handleChevronClick(false)}
          >
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
                d="M4.5 12.75l7.5-7.5 7.5 7.5m-15 6l7.5-7.5 7.5 7.5"
              />
            </svg>
          </button>
        </div>
      {/if}
    </div>

    <div>
      <button type="button" class="btn btn-primary invisible" aria-label="Settings">
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
            d="M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </button>
    </div>
  </div>

  <div class="mt-2">
    {#each Object.entries(logs) as [key] (key)}
      <div class={getTabClass(key)}>
        <button
          type="button"
          class="btn btn-ghost btn-xs normal-case"
          on:click={() => setActiveLog(key)}
        >
          {key.split(/[\\/]/).pop()}
        </button>

        <button class="btn btn-ghost btn-xs" on:click={() => removeLog(key)}>
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

  <div class="h-5/6 overflow-auto border-b border-base-300" id="lines">
    {#if activeLog}
      {#each logs[activeLog] as line}
        <p class="whitespace-nowrap font-mono">{line}</p>
      {/each}
    {/if}
  </div>
</main>
