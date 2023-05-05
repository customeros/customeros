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

// todo implement ssr

const client = new ApolloClient({
  ssrMode: typeof window === 'undefined',
  cache: new InMemoryCache({
    typePolicies: {
      Contact: {
        fields: {
          timelineEvents: {
            keyArgs: false,
            merge(existing = [], incoming) {
              return [...incoming, ...existing];
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
          dashboardView_Contacts: {
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
          dashboardView_Organizations: {
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
  assumeImmutableResults: true,
  connectToDevTools: true,
  credentials: 'include',
});

export default client;
