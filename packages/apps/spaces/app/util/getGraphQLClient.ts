import { GraphQLClient } from 'graphql-request';

export const getGraphQLClient = () => {
  return new GraphQLClient(
    `${process.env.NEXT_PUBLIC_SSR_PUBLIC_PATH}/customer-os-api/query`,
  );
};
