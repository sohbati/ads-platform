import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  use: {
    baseURL: "http://localhost:8100",
  },
  webServer: {
    command: "npm run dev -- --port 8100",
    url: "http://localhost:8100",
    reuseExistingServer: true,
  },
});
