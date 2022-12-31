<script lang="ts">
  import { toasts } from "../../stores/toastStore";
  import Toast from "./Toast.svelte";
</script>

{#if toasts}
  <section class="fixed right-4 bottom-0" style="min-width: 25%">
    {#each $toasts as toast (toast.id)}
      <Toast {toast} />
    {/each}
  </section>
{/if}
