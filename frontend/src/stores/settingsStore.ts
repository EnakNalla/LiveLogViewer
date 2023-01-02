import { writable } from "svelte/store";
import type { backend } from "../../wailsjs/go/models";

export const settings = writable<backend.Settings>({
  theme: "none",
  tailThreshold: 20,
  tailLines: 100,
  highlightErrors: true,
  highlightWarnings: false
});

export const setSettings = (s: backend.Settings) => settings.set(s);
