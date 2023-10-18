import { GraphQLClient } from 'graphql-request';
import { RequestConfig } from 'graphql-request/src/types';

// If request path will change and no longer match url we'd need to introduce variable
export const getGraphQLClient = (params?: RequestConfig) => {
  return new GraphQLClient(`/customer-os-api/query`, params);
};
