import { headers } from 'next/headers';

import { GraphQLClient } from 'graphql-request';

export const getServerGraphQLClient = () => {
  return new GraphQLClient(
    `${process.env.SSR_PUBLIC_PATH}/customer-os-api/query`,
    {
      credentials: 'include',
      headers: headers(),
    },
  );
};
