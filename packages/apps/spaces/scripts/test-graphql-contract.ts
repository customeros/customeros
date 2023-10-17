import type { CodegenConfig } from '@graphql-codegen/cli';
import path from 'path';
import dotenv from 'dotenv';

dotenv.config({
  path: path.join(__dirname, '../', '.env.development'),
});

const config: CodegenConfig = {
  overwrite: false,
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
    'scripts/temp/temp-graphql.ts': {
      plugins: ['typescript'],
    },
  },
  hooks: {
    afterAllFileWrite: ['prettier --write'],
  },
};
export default config;
