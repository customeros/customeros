import { GraphQLClient } from 'graphql-request';
import { headers } from 'next/headers';

export const getServerGraphQLClient = () => {
  return new GraphQLClient(
    `${process.env.SSR_PUBLIC_PATH}/customer-os-api/query`,
    {
      credentials: 'include',
      headers: headers(),
    },
  );
};
