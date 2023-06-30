import type { CodegenConfig } from '@graphql-codegen/cli';
import path from 'path';
import dotenv from 'dotenv';

dotenv.config({
  path: path.join(__dirname, '.env.development'),
});

const config: CodegenConfig = {
  overwrite: true,
  schema: [
    {
      [`${process.env.CUSTOMER_OS_API_PATH}/query`]: {
        headers: {
          'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY as string,
          'X-Openline-USERNAME': 'development@openline.ai',
        },
      },
    },
  ],
  documents: './app/**/*.graphql',
  generates: {
    'app/types/__generated__/graphql.types.ts': {
      plugins: ['typescript'],
    },
    './': {
      preset: 'near-operation-file',
      config: {
        namingConvention: {
          transformUnderscore: true,
        },
      },
      presetConfig: {
        baseTypesPath: 'app/types/__generated__/graphql.types.ts',
      },
      plugins: ['typescript-operations', 'typescript-react-query'],
    },
  },
  hooks: {
    afterOneFileWrite: ['prettier --write'],
  },
};
export default config;
