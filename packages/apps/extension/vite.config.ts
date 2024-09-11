import { defineConfig } from "vite";

export default defineConfig({
  build: {
    outDir: "dist",
    rollupOptions: {
      input: {
        background: "./src/background.ts",
        contentScript: "./src/contentScript.ts",
      },
      output: {
        entryFileNames: "[name].js",
        format: "es",
      },
    },
  },
  resolve: {
    alias: {
      "@": "/src",
    },
  },
});
