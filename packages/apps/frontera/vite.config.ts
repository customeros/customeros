import path from 'path';
import { cpus } from 'node:os';
import react from '@vitejs/plugin-react';
import { Plugin, defineConfig } from 'vite';
import graphqlLoader from 'vite-plugin-graphql-loader';

interface SourcemapExclude {
  excludeNodeModules?: boolean;
}

export function sourcemapExclude(opts?: SourcemapExclude): Plugin {
  return {
    name: 'sourcemap-exclude',
    transform(code: string, id: string) {
      if (opts?.excludeNodeModules && id.includes('node_modules')) {
        return {
          code,
          // https://github.com/rollup/rollup/blob/master/docs/plugin-development/index.md#source-code-transformations
          map: { mappings: '' },
        };
      }
    },
  };
}

// https://vitejs.dev/config/
export default defineConfig({
  // uncommenting the build underneath to produce sourcmaps will cause the build process to fail due to a bug in vite
  // check this issue: https://github.com/vitejs/vite/issues/2433
  build: {
    sourcemap: true,
    rollupOptions: {
      maxParallelFileOps: Math.max(1, cpus().length - 1),
      output: {
        manualChunks: (id) => {
          if (id.includes('node_modules')) {
            return 'vendor';
          }
        },
        sourcemapIgnoreList: (relativeSourcePath) => {
          const normalizedPath = path.normalize(relativeSourcePath);

          return normalizedPath.includes('node_modules');
        },
      },
      onwarn(warning, defaultHandler) {
        if (warning.code === 'SOURCEMAP_ERROR') {
          return;
        }

        defaultHandler(warning);
      },
    },
  },
  plugins: [
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
    sourcemapExclude({ excludeNodeModules: true }),
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
    },
  },
});
