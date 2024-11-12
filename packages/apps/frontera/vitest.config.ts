import path from 'path';
import { defineConfig } from 'vitest/config';
import graphqlLoader from 'vite-plugin-graphql-loader';

export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    name: 'Frontera',
    include: ['**/*.integration.ts'],
    testTimeout: 30000,
    setupFiles: ['src/store/vitest-hooks.ts'],
    hookTimeout: 60000,
  },
  plugins: [graphqlLoader()],
  resolve: {
    alias: {
      '@store': path.resolve(__dirname, './src/store'),
      '@graphql/types': path.resolve(
        __dirname,
        './src/routes/src/types/__generated__/graphql.types.ts',
      ),
    },
  },
});
