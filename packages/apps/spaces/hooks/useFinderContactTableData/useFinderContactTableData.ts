import {
  Contact,
  DashboardView_ContactsQueryVariables,
  useDashboardView_ContactsQuery,
} from './types';
import { ApolloError, NetworkStatus } from '@apollo/client';
import {
  Filter,
  InputMaybe,
  SortBy,
} from '../../graphQL/__generated__/generated';
import { useRecoilValue } from 'recoil';
import { finderContactTableSortingState } from '../../state/finderTables';


interface Result {
  data: Array<Contact> | null;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: {
    variables: DashboardView_ContactsQueryVariables;
  }) => void;
  refetchData: any;
  variables: DashboardView_ContactsQueryVariables;
  networkStatus?: NetworkStatus;
  totalElements: null | number;
}

export const useFinderContactTableData = (filters?: Filter[]): Result => {
  const initialVariables = {
    pagination: {
      page: 1,
      limit: 20,
    },
    where: undefined as InputMaybe<Filter> | undefined,
  };
  const sortingState = useRecoilValue(finderContactTableSortingState);
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
    useDashboardView_ContactsQuery({
      fetchPolicy: 'cache-and-network',
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
      data: data?.dashboardView_Contacts?.content || [],
      totalElements: data?.dashboardView_Contacts?.totalElements || null,
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
      totalElements: data?.dashboardView_Contacts?.totalElements || null,
    };
  }

  return {
    //@ts-expect-error revisit later, not matching generated types
    data: data?.dashboardView_Contacts?.content,
    totalElements: data?.dashboardView_Contacts?.totalElements || null,
    fetchMore,
    loading,
    error: null,
    variables: variables || initialVariables,
    refetchData: refetch,
    networkStatus,
  };
};
