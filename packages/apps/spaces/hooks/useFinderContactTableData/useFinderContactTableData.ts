import { ApolloError, NetworkStatus } from '@apollo/client';
import { Pagination } from './types';
import {
  Contact,
  GetContactListQueryVariables,
  useGetContactListQuery,
} from '../../graphQL/__generated__/generated';

interface Props {
  pagination: Pagination;
  searchTerm: string;
}

interface Result {
  data: Array<Contact> | null;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: { variables: GetContactListQueryVariables }) => void;
  variables: GetContactListQueryVariables;
  networkStatus?: NetworkStatus;
  totalElements: null | number;
}
export const useFinderContactTableData = ({
  pagination,
  searchTerm,
}: Props): Result => {
  const initialVariables = {
    pagination: {
      page: 0,
      limit: 10,
    },
    searchTerm,
  };
  const { data, loading, error, refetch, variables, fetchMore, networkStatus } =
    useGetContactListQuery({
      fetchPolicy: 'cache-first',
      notifyOnNetworkStatusChange: true,
      variables: initialVariables,
    });

  if (loading) {
    return {
      loading: true,
      error: null,
      //@ts-expect-error revisit later, not matching generated types
      data: data?.contacts?.content || [],
      totalElements: data?.contacts?.totalElements || null,
      fetchMore,
      variables: variables || initialVariables,
      networkStatus,
    };
  }

  if (error) {
    return {
      error,
      loading: false,
      variables: variables || initialVariables,
      networkStatus,
      data: null,
      fetchMore,
      totalElements: data?.contacts?.totalElements || null,
    };
  }

  return {
    //@ts-expect-error revisit later, not matching generated types
    data: data?.contacts?.content,
    totalElements: data?.contacts?.totalElements || null,
    fetchMore,
    loading,
    error: null,
    variables: variables || initialVariables,
    refetch,
    networkStatus,
  };
};
