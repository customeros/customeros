import { defineConfig } from "vite";
import { resolve } from 'path';
import fs from 'fs';

export default defineConfig({
  build: {
    outDir: "dist",
    rollupOptions: {
      input: {
        background: "./src/background.ts",
        contentScript: "./src/contentScript.ts",
        sidepanel: "./src/sidepanel.js",
        sidepanelHtml: "./src/sidepanel.html", // Add this line
      },
      output: {
        entryFileNames: "[name].js",
        assetFileNames: "[name].[ext]",
        format: "es",
      },
    },
  },
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  plugins: [
    {
      name: 'html-transform',
      transformIndexHtml: {
        enforce: 'pre',
        transform(html, ctx) {
          if (ctx.path.endsWith('sidepanel.html')) {
            return {
              html,
              tags: [
                {
                  tag: 'script',
                  attrs: { src: 'sidepanel.js', type: 'module' },
                  injectTo: 'body',
                },
              ],
            };
          }
        },
      },
    },
    {
      name: 'copy-sidepanel-html',
      writeBundle() {
        fs.copyFileSync('./src/sidepanel.html', './dist/sidepanel.html');
      },
    },
  ],
});
