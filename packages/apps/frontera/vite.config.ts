import path from 'path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import graphqlLoader from 'vite-plugin-graphql-loader';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        parserOpts: {
          plugins: ['decorators'],
        },
      },
    }),
    graphqlLoader(),
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
      '@organizations': path.resolve(
        __dirname,
        './src/routes/organizations/src',
      ),
      '@organization': path.resolve(__dirname, './src/routes/organization/src'),
      '@renewals': path.resolve(__dirname, './src/routes/renewals/src'),
      '@settings': path.resolve(__dirname, './src/routes/settings/src'),
      '@utils': path.resolve(__dirname, './src/utils'),
      '@invoices': path.resolve(__dirname, './src/routes/invoices/src'),
      '@opportunities': path.resolve(__dirname, './src/routes/prospects/src'),
    },
  },
});
