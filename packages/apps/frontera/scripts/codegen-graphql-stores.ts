import type { CodegenConfig } from '@graphql-codegen/cli';

import path from 'path';
import dotenv from 'dotenv';

dotenv.config({
  path: path.join(__dirname, '../', '.env'),
});

const config: CodegenConfig = {
  overwrite: true,
  schema: [
    {
      [`${process.env.CUSTOMER_OS_API_PATH}/query`]: {
        headers: {
          'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY as string,
          'X-Openline-USERNAME': 'edi@customeros.ai',
        },
      },
    },
  ],
  documents: './src/store/**/*.graphql',
  generates: {
    'src/routes/src/types/__generated__/graphql.types.ts': {
      plugins: ['typescript'],
    },
    './': {
      preset: 'near-operation-file',
      config: {
        namingConvention: {
          transformUnderscore: true,
        },
        exposeDocument: true,
      },
      presetConfig: {
        baseTypesPath: 'src/routes/src/types/__generated__/graphql.types.ts',
      },
      plugins: ['typescript-operations'],
    },
  },
  hooks: {
    afterOneFileWrite: ['prettier --write', 'echo'],
  },
};
export default config;
