import type { Settings } from "$lib/types/settings";

export const DEFAULT_SETTINGS: Settings = {
  autostart: false,
  log_level: "INFO",
  server: {
    host: "localhost",
    port: 8080,
  },
};
