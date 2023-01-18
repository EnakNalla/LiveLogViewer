<script lang="ts">
  import { RemoveLog, SelectLog } from "../wailsjs/go/backend/App";
  import { EventsOn } from "../wailsjs/runtime";
  import Settings from "./components/Settings.svelte";
  import ToastContainer from "./components/Toast/ToastContainer.svelte";
  import { settings } from "./stores/settingsStore";
  import { addToast } from "./stores/toastStore";

  // TODO: very specific maybe regex?
  const levels = ["[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[TRACE]", "[FATAL]"];

  let logs: { [key: string]: string[] } = {};

  let activeLog = "";
  let paused = false;

  let searchTerm = "";
  let searchResults = 0;
  let searchIndex = 0;

  const togglePaused = () => (paused = !paused);

  EventsOn("error", err => addToast(err, "alert-error"));

  const openLog = async () => {
    const response = await SelectLog();
    if (!response.success) {
      addToast(response.error, "alert-error");
      return;
    }

    const path = response.data.path;
    logs[path] = response.data.lines;

    setTimeout(() => setActiveLog(path), 50);
    setTimeout(() => scrollToBottom(path), 200);

    EventsOn(path, line => {
      const log = logs[path];
      logs[path] = [...log, line];

      if (!paused && activeLog === path) scrollToBottom(path);
    });
  };

  const scrollToBottom = (id: string) => {
    const el = document.getElementById(id);
    el.scrollTo(0, el.scrollHeight);
  };

  const setActiveLog = (key: string) => {
    clearSearch();
    activeLog = key;
    logs = { ...logs };

    document.querySelector(".visible")?.classList.replace("visible", "hidden");
    document.getElementById(key).classList.replace("hidden", "visible");
  };

  const removeLog = (path: string) => {
    RemoveLog(path);

    delete logs[path];

    logs = { ...logs };

    if (path === activeLog) {
      const entries = Object.entries(logs);
      activeLog = entries.length >= 1 ? entries[0][0] : "";
      if (activeLog) setActiveLog(activeLog);
    }
  };

  const getTabClass = (tab: string) =>
    tab === activeLog
      ? "tab tab-bordered tab-lifted tab-active !bg-base-300"
      : "tab tab-bordered tab-lifted";

  const replaceInvariant = (str: string, term: string) => {
    var esc = str.replace(/[-\/\\^$*+?.()|[\]{}]/, "\\$&");
    var reg = new RegExp(esc, "i");
  };

  const handleSearch = () => {
    if (!searchTerm) {
      clearSearch();
      return;
    }

    paused = true;
    let results = 0;
    const lines = document.querySelectorAll("p");

    for (let i = 0; i < lines.length; i++) {
      let line = lines[i];
      let condition = $settings.ignoreCase
        ? line.innerText.toLowerCase().includes(searchTerm.toLowerCase())
        : line.innerText.includes(searchTerm);

      if (condition) {
        var esc = searchTerm.replace(/[-\/\\^$*+?.()|[\]{}]/, "\\$&");
        var reg = new RegExp(esc, "gi");
        line.innerHTML = line.innerText.replace(
          reg,
          `<span class="text-black bg-yellow-400" data-index="${results}">$&</span>`
        );
        results++;
      }
    }

    searchResults = results - 1;

    if (searchResults <= 0) {
      addToast("No results!", "alert-warning");
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
    } else if (searchIndex > searchResults) {
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

  const parseLine = (line: string) => {
    const level = levels.find(level => line.includes(level))!;
    if (!level) return line;

    let colour = "";
    if (line.includes("ERROR") || line.includes("FATAL")) colour = "text-red-500";
    if (line.includes("WARN")) colour = "text-yellow-500";
    if (line.includes("INFO")) colour = "text-blue-500";
    if (line.includes("DEBUG")) colour = "text-green-500";
    if (line.includes("TRACE")) colour = "text-purple-500";

    return line.replace(level, `<span class='${colour}'>${level}</span>`);
  };
</script>

<main class="mx-4 h-full">
  <div class="flex justify-between pt-2">
    <div>
      <button type="button" class="btn-primary btn" aria-label="Open log" on:click={openLog}>
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
      <Settings />
    </div>
  </div>

  <div class="mt-2">
    <div class="mt-2" />
    {#each Object.entries(logs) as [key]}
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

  {#each Object.entries(logs) as [key, lines] (key)}
    <div
      class="rounded-b-box rounded-tr-box bg-base-300 hidden h-5/6 overflow-auto p-2 pb-8"
      id={key}
    >
      {#each lines as line}
        {#if $settings.highlightLevels}
          <p>{@html parseLine(line)}</p>
        {:else}
          <p>{line}</p>
        {/if}
      {/each}
    </div>
  {/each}

  <ToastContainer />
</main>
