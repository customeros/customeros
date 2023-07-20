import { GraphQLClient } from 'graphql-request';

// If request path will change and no longer match url we'd need to introduce variable
export const getGraphQLClient = () => {
  return new GraphQLClient(`/customer-os-api/query`);
};
