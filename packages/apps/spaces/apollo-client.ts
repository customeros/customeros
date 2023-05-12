import { ApolloClient, HttpLink, InMemoryCache, from } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';

export const httpLink = new HttpLink({
  uri: `/customer-os-api/query`,
  fetchOptions: {
    credentials: 'include',
  },
});

export const authLink = setContext((_, { headers }) => {
  return {
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
  };
});

const client = new ApolloClient({
  ssrMode: typeof window === 'undefined',
  cache: new InMemoryCache({
    typePolicies: {
      Contact: {
        fields: {
          timelineEvents: {
            keyArgs: ['id'],
            merge(
              existing = [],
              incoming=[],
            ) {

              const merged = existing ? existing.slice(0) : [];
              const existingIds = existing ? existing.map(item => item.__ref) : [];
              incoming.forEach((item) => {
                if (existingIds.indexOf(item.__ref) < 0) {
                  merged.push(item);
                }
              });
              console.log('🏷️ ----- merged: '
                  , merged);
              return merged;
            },
          },
        },
      },

      Organization: {
        fields: {
          timelineEvents: {
            keyArgs: false,
            merge(existing = [], incoming) {
              return [...incoming, ...existing];
            },
          },
        },
      },
      Query: {
        fields: {
          dashboardView: {
            keyArgs: false,
            merge(
              existing = { content: [] },
              incoming,
              {
                args: {
                  // @ts-expect-error look into it later
                  pagination: { page, limit },
                },
              },
            ) {
              if (page === 1) return incoming;
              return {
                ...existing,
                content: [...existing.content, ...incoming.content],
              };
            },
          },
        },
      },
    },
  }),
  link: from([authLink, httpLink]),
  queryDeduplication: true,
  assumeImmutableResults: false,
  connectToDevTools: true,
  credentials: 'include',
});

export default client;
