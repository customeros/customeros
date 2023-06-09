import { ApolloError, NetworkStatus } from '@apollo/client';
import {
  DashboardView_OrganizationsQueryVariables,
  Organization,
  useDashboardView_OrganizationsQuery,
} from './types';
import {
  Filter,
  InputMaybe,
  SortBy,
} from '../../graphQL/__generated__/generated';
import { useRecoilValue } from 'recoil';
import { finderOrganizationTableSortingState } from '../../state/finderTables';

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

export const useFinderOrganizationTableData = (filters?: Filter[]): Result => {
  const initialVariables = {
    pagination: {
      page: 1,
      limit: 20,
    },
    where: undefined as InputMaybe<Filter> | undefined,
  };
  const sortingState = useRecoilValue(finderOrganizationTableSortingState);
  const sortBy: SortBy | undefined = sortingState.column
    ? {
        by: sortingState.column,
        direction: sortingState.direction,
        caseSensitive: false,
      }
    : undefined;
  if (filters && filters.length > 0) {
    initialVariables.where = { AND: filters } as Filter;
  }
  const { data, loading, error, refetch, variables, fetchMore, networkStatus } =
    useDashboardView_OrganizationsQuery({
      fetchPolicy: 'network-only',
      notifyOnNetworkStatusChange: true,
      variables: {
        pagination: initialVariables.pagination,
        where: initialVariables.where,
        sort: sortBy,
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
