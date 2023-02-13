import type { CodegenConfig } from '@graphql-codegen/cli';
const path = require('path');
import dotenv from 'dotenv';

dotenv.config({
  path: path.join(__dirname, '.env.dev'),
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
  documents: 'src/**/*.tsx',
  generates: {
    'src/graphQL/graphql.schema.json': {
      plugins: ['introspection'],
    },
    'src/graphQL/types.ts': {
      plugins: ['typescript'],
    },
    'src/graphQL/hooks.ts': {
      plugins: ['typescript-operations', 'typescript-react-apollo'],
      config: { withHooks: true },
    },
    'src/': {
      preset: 'near-operation-file',
      presetConfig: { extension: '.generated.tsx', baseTypesPath: 'types.ts' },
      plugins: ['typescript-operations', 'typescript-react-apollo'],
    },
  },
};
export default config;
