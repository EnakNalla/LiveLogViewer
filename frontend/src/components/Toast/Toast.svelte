<script lang="ts">
  import { type Toast, removeToast } from "../../stores/toastStore";

  export let toast: Toast;
</script>

<div class={`alert ${toast.theme} m-2 block w-full p-0.5`}>
  {#if toast.title}
    <div class="flex w-full justify-between">
      <h3 class="ml-2 text-lg font-bold">{toast.title}</h3>
      <button
        type="button"
        class="btn-ghost btn-xs btn"
        aria-label="Dismiss"
        on:click={() => removeToast(toast.id)}
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
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {:else}
    <div class="flex w-full justify-end">
      <button
        type="button"
        class="btn-ghost btn-xs btn"
        aria-label="Dismiss"
        on:click={() => removeToast(toast.id)}
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
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {/if}

  <p class="ml-2">{toast.message}</p>
</div>
