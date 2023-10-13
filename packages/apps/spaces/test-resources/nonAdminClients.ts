import { GraphQLClient } from 'graphql-request';

export const headers = {
  'X-Openline-API-KEY': 'dd9d2474-b4a9-4799-b96f-73cd0a2917e4',
  'X-Openline-USERNAME': 'silviu@openline.ai',
};

export const authenticatedClient = new GraphQLClient(
  'http://127.0.0.1:10000/query',
  {
    credentials: 'include',
    headers,
  },
);

export const unauthenticatedClient = new GraphQLClient(
  'http://127.0.0.1:10000/query',
  {
    credentials: 'include',
  },
);
