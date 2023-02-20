import type { CodegenConfig } from '@graphql-codegen/cli';
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
  documents: 'src/**/**/graphQL/**',
  generates: {
    'src/graphQL/graphql.schema.json': {
      plugins: ['introspection'],
    },
    'src/graphQL/generated.ts': {
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
