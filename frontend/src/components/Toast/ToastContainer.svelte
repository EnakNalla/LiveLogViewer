<script lang="ts">
  import { toasts } from "../../stores/toastStore";
  import Toast from "./Toast.svelte";
</script>

{#if toasts}
  <section class="fixed right-4 bottom-0 w-1/4">
    {#each $toasts as toast (toast.id)}
      <Toast {toast} />
    {/each}
  </section>
{/if}
