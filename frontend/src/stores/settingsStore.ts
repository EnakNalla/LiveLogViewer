import { writable } from "svelte/store";
import type { backend } from "../../wailsjs/go/models";

export const DEFAULT_SETTINGS: backend.Settings = {
  theme: "none",
  tailThreshold: 5,
  tailLines: 100,
  highlightLevels: true,
  pollingEnabled: true,
  pollInterval: 1000,
  ignoreCase: true,
  textWrap: false,
  lineNumbers: true
};

export const settings = writable<backend.Settings>(DEFAULT_SETTINGS);

export const setSettings = (s: backend.Settings) => settings.set(s);
