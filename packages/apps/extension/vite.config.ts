import path from "path";

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";
import fs from "fs";

export default defineConfig({
  build: {
    outDir: "dist",
    rollupOptions: {
      input: {
        background: resolve(__dirname, "src/background.ts"),
        contentScript: resolve(__dirname, "src/contentScript.ts"),
        main: resolve(__dirname, "src/sidepanel/main.tsx"),
      },
      output: {
        entryFileNames: (chunkInfo) => {
          if (chunkInfo.name === "main") {
            return "sidepanel/[name].js";
          }
          return "[name].js";
        },
        chunkFileNames: "sidepanel/chunks/[name].js",
        assetFileNames: (assetInfo) => {
          const info = assetInfo.name?.split(".") ?? [];
          const extType = info[info.length - 1];

          if (extType === "css") {
            return "sidepanel/styles/tailwind.[ext]";
          }
          return "assets/[name].[ext]";
        },
      },
    },
  },
  css: {
    postcss: "./postcss.config.js",
  },
  plugins: [
    react(),
    {
      name: "copy-sidepanel-html",
      writeBundle() {
        // Copy HTML file
        let html = fs.readFileSync("./src/sidepanel/index.html", "utf-8");
        fs.mkdirSync("./dist/sidepanel", { recursive: true });
        fs.writeFileSync("./dist/sidepanel/index.html", html);

        // Copy assets
        if (fs.existsSync("./src/sidepanel/assets")) {
          fs.mkdirSync("./dist/src/assets", { recursive: true });
          fs.cpSync("./src/sidepanel/assets", "./dist/src/assets", {
            recursive: true,
          });
        }

        // Copy UI components
        if (fs.existsSync("./src/sidepanel/ui")) {
          fs.mkdirSync("./dist/sidepanel/ui", { recursive: true });
          fs.cpSync("./src/sidepanel/ui", "./dist/sidepanel/ui", {
            recursive: true,
          });
        }
      },
    },
  ],
  resolve: {
    alias: {
      "@ui": path.resolve(__dirname, "./src/sidepanel/ui"),
    },
  },
});
