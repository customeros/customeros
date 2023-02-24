import { ApolloClient, HttpLink, InMemoryCache } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';

const httpLink = new HttpLink({
  uri: `/customer-os-api/query`,
  fetchOptions: {
    credentials: 'include',
  },
});

const authLink = setContext((_, { headers }) => {
  return {
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
  };
});

const client = new ApolloClient({
  cache: new InMemoryCache({
    typePolicies: {
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
                  // @ts-expect-error look into it later
                  searchTerm,
                },
              },
            ) {
              if (page === 0) return incoming;
              if (searchTerm && page === 0) {
                const listLength = incoming.totalElements;
                if (listLength <= limit) {
                  return incoming;
                }
              }
              if (searchTerm && page > 0) {
                return {
                  ...existing,
                  content: [...existing.content, ...incoming.content],
                };
              }

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
  link: authLink.concat(httpLink),
  queryDeduplication: true,
  assumeImmutableResults: true,
  connectToDevTools: true,
  credentials: 'include',
});

export default client;
