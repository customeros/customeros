import { GraphQLClient } from 'graphql-request';
// @ts-expect-error - looks like types are causing errors after bumping version
import { RequestConfig } from 'graphql-request/src/types';

// If request path will change and no longer match url we'd need to introduce variable
export const getGraphQLClient = (params?: RequestConfig) => {
  return new GraphQLClient(
    `${import.meta.env.VITE_MIDDLEWARE_API_URL}/customer-os-api`,
    {
      ...params,
      headers: {
        ...params?.headers,
        // 'X-Frontera-Auth': 'true',
        Authorization: `Bearer ${window?.__COS_SESSION__?.sessionToken}`,
        'X-Openline-USERNAME': window?.__COS_SESSION__?.email,
      },
    },
  );
};
