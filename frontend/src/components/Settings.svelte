<script lang="ts">
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
      <h3 class="text-lg font-bold">Settings</h3>
      <label for="settings-modal" aria-label="Close settings" class="btn-ghost btn">
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
  </label>
</label>
