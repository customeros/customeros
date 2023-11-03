import { GraphQLClient } from 'graphql-request';

const x = 'y';

export const headers = {
  'X-Openline-API-KEY': 'dd9d2474-b4a9-4799-b96f-73cd0a2917e4',
  'X-Openline-TENANT': 'openline',
};
export const authenticated_client = new GraphQLClient(
  'http://127.0.0.1:10000/query',
  {
    credentials: 'include',
    headers,
  },
);

export const unauthenticated_client = new GraphQLClient(
  'http://127.0.0.1:10000/query',
  {
    credentials: 'include',
  },
);
