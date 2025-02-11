import path from 'path';
import fg from 'fast-glob';
import { cpus } from 'node:os';
import svgr from 'vite-plugin-svgr';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import graphqlLoader from 'vite-plugin-graphql-loader';

export default defineConfig(({ mode }) => ({
  build: {
    sourcemap: mode === 'production',
    rollupOptions: {
      maxParallelFileOps: Math.max(1, cpus().length - 1),
      output: {
        manualChunks: (id) => {
          if (id.includes('node_modules')) {
            return 'vendor';
          }
        },
        sourcemapExcludeSources: true,
        sourcemapIgnoreList: (relativePath) => {
          if (relativePath.includes('node_modules')) {
            return true;
          }

          return false;
        },
      },
    },
  },
  server: {
    allowedHosts: ['localhost', '127.0.0.1', 'app.customeros.local'],
  },
  plugins: [
    {
      name: 'watch-external',
      async buildStart() {
        const files = await fg(['public/**/*']);

        for (const file of files) {
          this.addWatchFile(file);
        }
      },
    },
    react({
      babel: {
        plugins: [
          [
            '@babel/plugin-proposal-decorators',
            {
              version: '2023-05',
            },
          ],
        ],
      },
    }),
    graphqlLoader(),
    svgr(),
  ],
  resolve: {
    alias: {
      '@ui': path.resolve(__dirname, './src/ui'),
      '@shared': path.resolve(__dirname, './src/routes/src'),
      '@assets': path.resolve(__dirname, './src/assets'),
      '@store': path.resolve(__dirname, './src/store'),
      '@graphql/types': path.resolve(
        __dirname,
        './src/routes/src/types/__generated__/graphql.types.ts',
      ),
      '@finder': path.resolve(__dirname, './src/routes/finder/src'),
      '@organization': path.resolve(__dirname, './src/routes/organization/src'),
      '@renewals': path.resolve(__dirname, './src/routes/renewals/src'),
      '@settings': path.resolve(__dirname, './src/routes/settings/src'),
      '@utils': path.resolve(__dirname, './src/utils'),
      '@invoices': path.resolve(__dirname, './src/routes/invoices/src'),
      '@opportunities': path.resolve(__dirname, './src/routes/prospects/src'),
      '@domain': path.resolve(__dirname, './src/domain'),
      '@infra': path.resolve(__dirname, './src/infra'),
    },
  },
}));
