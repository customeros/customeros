import { ApolloError, NetworkStatus } from '@apollo/client';
import {
  DashboardView_OrganizationsQueryVariables,
  Organization,
  useDashboardView_OrganizationsQuery,
} from './types';
import { Filter, InputMaybe, SortBy } from '@spaces/graphql';

interface Result {
  data: Array<Organization> | null;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: {
    variables: DashboardView_OrganizationsQueryVariables;
  }) => void;
  refetchData: any;
  variables: DashboardView_OrganizationsQueryVariables;
  networkStatus?: NetworkStatus;
  totalElements: null | number;
}

export const useFinderOrganizationTableData = (
  filters?: Filter[],
  sortBy?: SortBy,
): Result => {
  const initialVariables = {
    pagination: {
      page: 1,
      limit: 40,
    },
    where: undefined as InputMaybe<Filter> | undefined,
  };
  const { data, loading, error, refetch, variables, fetchMore, networkStatus } =
    useDashboardView_OrganizationsQuery({
      fetchPolicy: 'network-only',
      notifyOnNetworkStatusChange: true,
      variables: {
        pagination: initialVariables.pagination,
        where: !filters?.length
          ? (initialVariables.where as Filter)
          : ({
              AND: filters,
            } as Filter),
        sort: sortBy ?? undefined,
      },
    });

  if (loading) {
    return {
      loading: true,
      error: null,
      //@ts-expect-error revisit later, not matching generated types
      data: data?.dashboardView_Organizations?.content || [],
      totalElements: data?.dashboardView_Organizations?.totalElements || null,
      fetchMore,
      refetchData: refetch,
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
      refetchData: refetch,
      totalElements: data?.dashboardView_Organizations?.totalElements || null,
    };
  }

  return {
    //@ts-expect-error revisit later, not matching generated types
    data: data?.dashboardView_Organizations?.content,
    totalElements: data?.dashboardView_Organizations?.totalElements || null,
    fetchMore,
    loading,
    error: null,
    variables: variables || initialVariables,
    refetchData: refetch,
    networkStatus,
  };
};
