<script lang="ts">
  import { GetSettings, WriteSettings } from "../../wailsjs/go/main/App";
  import { onMount } from "svelte";
  import { addToast } from "../stores/toastStore";
  import type { main } from "../../wailsjs/go/models";

  let settings: main.Settings;
  const htmlEl = document.querySelector("html");

  onMount(async () => {
    const result = await GetSettings();
    if (!result.success) {
      addToast("Failed to load config file", "alert-error");
      return;
    }

    settings = result.data;

    if (settings.theme === "none") {
      settings.theme = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
      WriteSettings(settings);
      htmlEl.dataset.theme = settings.theme;
    } else {
      htmlEl.dataset.theme = settings.theme;
    }
  });

  const toggleTheme = async () => {
    settings.theme = settings.theme === "dark" ? "light" : "dark";
    const result = await WriteSettings(settings);
    if (!result.success) {
      addToast("Failed to save settings", "alert-error");
    } else {
      htmlEl.dataset.theme = settings.theme;
    }
  };
</script>

<label for="settings-modal" class="btn-primary btn" aria-label="Open settings">
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
</label>

<input type="checkbox" id="settings-modal" class="modal-toggle" />

<label for="settings-modal" class="modal cursor-pointer">
  <label class="modal-box relative" for="">
    <div class="flex justify-between">
      <h3 class="align-middle text-lg font-bold">Settings</h3>
      <label for="settings-modal" aria-label="Close settings" class="btn-ghost btn-sm btn">
        <svg
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="h-6 w-6"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </label>
    </div>

    <div class="m-2">
      <label class="label cursor-pointer">
        <span class="label-text">
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
              d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"
            />
          </svg>
        </span>
        <input
          type="checkbox"
          class="toggle"
          checked={settings?.theme !== "dark"}
          on:click={toggleTheme}
        />
        <span class="label-text-alt">
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
              d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
            />
          </svg>
        </span>
      </label>
    </div>

    <div class="modal-action">
      <label for="settings-modal" class="btn">Close</label>
    </div>
  </label>
</label>
