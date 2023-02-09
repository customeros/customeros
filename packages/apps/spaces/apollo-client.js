import { ApolloClient, InMemoryCache } from '@apollo/client';

const cache = new InMemoryCache({
  // todo
});

const client = new ApolloClient({
  cache,
  uri: process.env.CUSTOMER_OS_API_PATH,
  queryDeduplication: true,
  assumeImmutableResults: true,
  connectToDevTools: true,
  query: {
    errorPolicy: 'all',
  },
});

export default client