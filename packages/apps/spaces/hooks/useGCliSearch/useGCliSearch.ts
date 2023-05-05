import { ApolloError, NetworkStatus } from 'apollo-client';
import {
  GCliSearchQueryVariables,
  GCliSearchResultItem,
  useGCliSearchQuery,
} from './types';

interface Props {
  searchTerm: string;
}

interface Result {
  data: Array<GCliSearchResultItem> | null;
  loading: boolean;
  error: ApolloError | null;
  variables: GCliSearchQueryVariables;
  networkStatus?: NetworkStatus;
  refetch?: (variables?: GCliSearchQueryVariables) => Promise<any>;
}

export const useGCliSearch = (): Result => {
  const initialVariables = {
    limit: 5,
    keyword: '',
  };
  const { data, loading, error, refetch, variables, networkStatus } =
    useGCliSearchQuery({
      fetchPolicy: 'cache-first',
      notifyOnNetworkStatusChange: true,
      variables: { limit: 5, keyword: '' },
    });

  if (loading) {
    return {
      loading: true,
      error: null,
      data: data?.gcli_Search || [],
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
    //@ts-expect-error revisit later, not matching generated types
    data: data?.gcli_Search,
    loading,
    error: null,
    variables: variables || initialVariables,
    refetch,
    networkStatus,
  };
};
