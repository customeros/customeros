import { ApolloError, NetworkStatus } from 'apollo-client';
import {
  useGetDashboardDataQuery,
  DashboardViewItem,
  GetDashboardDataQueryVariables,
  Pagination,
} from './types';

interface Props {
  pagination: Pagination;
  searchTerm: string;
}

interface Result {
  data: Array<DashboardViewItem> | null;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: { variables: GetDashboardDataQueryVariables }) => void;
  variables: GetDashboardDataQueryVariables;
  networkStatus?: NetworkStatus;
  totalElements: null | number;
}
export const useFinderTableData = ({
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
    useGetDashboardDataQuery({
      fetchPolicy: 'cache-first',
      notifyOnNetworkStatusChange: true,
      variables: initialVariables,
    });

  if (loading) {
    return {
      loading: true,
      error: null,
      //@ts-expect-error revisit later, not matching generated types
      data: data?.dashboardView?.content || [],
      totalElements: data?.dashboardView?.totalElements || null,
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
      totalElements: data?.dashboardView?.totalElements || null,
    };
  }

  return {
    //@ts-expect-error revisit later, not matching generated types
    data: data?.dashboardView?.content,
    totalElements: data?.dashboardView?.totalElements || null,
    fetchMore,
    loading,
    error: null,
    variables: variables || initialVariables,
    refetch,
    networkStatus,
  };
};
