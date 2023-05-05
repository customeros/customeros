import { ApolloError, NetworkStatus } from '@apollo/client';
import {
  DashboardView_OrganizationsQueryVariables,
  Organization,
  useDashboardView_OrganizationsQuery,
} from './types';

interface Result {
  data: Array<Organization> | null;
  loading: boolean;
  error: ApolloError | null;
  fetchMore: (data: {
    variables: DashboardView_OrganizationsQueryVariables;
  }) => void;
  variables: DashboardView_OrganizationsQueryVariables;
  networkStatus?: NetworkStatus;
  totalElements: null | number;
}

export const useFinderOrganizationTableData = (): Result => {
  const initialVariables = {
    pagination: {
      page: 1,
      limit: 60,
    },
  };
  const { data, loading, error, refetch, variables, fetchMore, networkStatus } =
    useDashboardView_OrganizationsQuery({
      fetchPolicy: 'cache-first',
      notifyOnNetworkStatusChange: true,
      variables: { pagination: initialVariables.pagination },
    });

  if (loading) {
    return {
      loading: true,
      error: null,
      //@ts-expect-error revisit later, not matching generated types
      data: data?.organizations?.content || [],
      totalElements: data?.dashboardView_Organizations?.totalElements || null,
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
    refetch,
    networkStatus,
  };
};
