import type { CodegenConfig } from '@graphql-codegen/cli';

import path from 'path';
import dotenv from 'dotenv';

dotenv.config({
  path: path.join(__dirname, '../', '.env.development'),
});

const config: CodegenConfig = {
  overwrite: true,
  schema: [
    {
      [`${process.env.CUSTOMER_OS_API_PATH}/query`]: {
        headers: {
          'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY as string,
          'X-Openline-USERNAME': 'edi@openline.ai',
        },
      },
    },
  ],
  documents: './app/**/*.graphql',
  generates: {
    'app/src/types/__generated__/graphql.types.ts': {
      plugins: ['typescript'],
    },
    './': {
      preset: 'near-operation-file',
      config: {
        namingConvention: {
          transformUnderscore: true,
        },
        exposeDocument: true,
        exposeQueryKeys: true,
        exposeMutationKeys: true,
        exposeFetcher: true,
        addInfiniteQuery: true,
        fetcher: 'graphql-request',
        reactQueryVersion: 5,
      },
      presetConfig: {
        baseTypesPath: 'app/src/types/__generated__/graphql.types.ts',
      },
      plugins: [
        'typescript-operations',
        'typescript-react-query',
        {
          add: {
            content:
              '// @ts-nocheck remove this when typscript-react-query plugin is fixed',
          },
        },
        'scripts/expose-cache-mutator.js',
      ],
    },
  },
  hooks: {
    afterOneFileWrite: ['prettier --write', 'echo'],
  },
};
export default config;
