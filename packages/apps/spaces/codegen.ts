import type { CodegenConfig } from '@graphql-codegen/cli';
// eslint-disable-next-line @typescript-eslint/no-var-requires
const path = require('path');
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
  documents: '**/**/graphQL/**/**',
  generates: {
    'graphQL/__generated__/graphql.schema.json': {
      plugins: ['introspection'],
    },
    'graphQL/__generated__/generated.ts': {
      plugins: [
        'typescript',
        'typescript-operations',
        'typescript-react-apollo',
      ],
      config: { withHooks: true },
    },
  },
};
export default config;
