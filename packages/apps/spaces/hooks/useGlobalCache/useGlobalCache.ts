import { Global_CacheQueryVariables, useGlobal_CacheQuery } from './types';
import { ApolloError, NetworkStatus } from '@apollo/client';

interface Props {
  searchTerm: string;
}

interface Result {
  data: any | null;
  loading: boolean;
  error: ApolloError | null;
  variables: Global_CacheQueryVariables;
  networkStatus?: NetworkStatus;
  refetch?: (variables?: Global_CacheQueryVariables) => Promise<any>;
}

export const useGlobalCache = (): Result => {
  const initialVariables = {};
  const { data, loading, error, variables, refetch, networkStatus } =
    useGlobal_CacheQuery({
      fetchPolicy: 'network-only',
      notifyOnNetworkStatusChange: true,
    });

  if (loading) {
    return {
      loading: true,
      error: null,
      data: [],
      variables: variables || initialVariables,
      refetch,
      networkStatus,
    };
  }

  if (error) {
    return {
      error,
      loading: false,
      variables: variables || initialVariables,
      networkStatus,
      refetch,
      data: null,
    };
  }

  return {
    data: data?.global_Cache,
    loading,
    error: null,
    variables: variables || initialVariables,
    refetch,
    networkStatus,
  };
};
