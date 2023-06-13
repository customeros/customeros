import { Global_CacheQueryVariables, useGlobal_CacheQuery } from './types';
import { ApolloError, NetworkStatus } from '@apollo/client';
import {
  GetContactListQueryVariables,
  useGetContactListLazyQuery,
  useGlobal_CacheLazyQuery,
} from '@spaces/graphql';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadGlobalCache: () => Promise<any>;
}

export const useGlobalCache = (): Result => {
  const [onLoadGlobalCache, { data, loading, error }] =
    useGlobal_CacheLazyQuery({
      fetchPolicy: 'network-only',
    });

  return {
    onLoadGlobalCache,
    loading,
    error,
  };
};
