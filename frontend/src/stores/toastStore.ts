import { writable } from "svelte/store";

type ToastTheme = "alert-error" | "alert-info" | "alert-warning" | "alert-success";
export interface Toast {
  id: number;
  title?: string;
  message: string;
  theme?: ToastTheme;
}

export const toasts = writable<Toast[]>([]);

export const addToast = (message: string, theme?: ToastTheme, title?: string) => {
  const id = Math.floor(Math.random() * 10000);

  if (!title) {
    switch (theme) {
      case "alert-error":
        title = "Error!";
        break;
      case "alert-info":
        title = "Info";
        break;
      case "alert-warning":
        title = "Warning";
        break;
      case "alert-success":
        title = "Success!";
        break;
      default:
        title = "Alert";
    }
  }

  toasts.update(all => [{ id, message, theme, title }, ...all]);

  setTimeout(() => removeToast(id), 2000);
};

export const removeToast = (id: number) => {
  toasts.update(all => all.filter(t => t.id !== id));
};
