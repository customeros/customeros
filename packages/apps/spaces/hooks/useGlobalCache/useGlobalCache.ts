import { ApolloError } from '@apollo/client';
import { useGlobal_CacheLazyQuery } from '@spaces/graphql';

import ApolloClient from '../../apollo-client';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadGlobalCache: () => Promise<any>;
}

export const useGlobalCache = (): Result => {
  const [onLoadGlobalCache, { data, loading, error }] =
    useGlobal_CacheLazyQuery({
      fetchPolicy: 'network-only',
      client: ApolloClient,
    });

  return {
    onLoadGlobalCache,
    loading,
    error,
  };
};
