import { ApolloError } from 'apollo-client';
import { GetTenantNameQuery, useGetTenantNameQuery } from './types';

interface Result {
  data: GetTenantNameQuery['tenant'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useTenantName = (): Result => {
  const { data, loading, error } = useGetTenantNameQuery();

  if (loading) {
    return {
      loading: true,
      error: null,
      data: null,
    };
  }

  if (error) {
    return {
      error,
      loading: false,
      data: null,
    };
  }

  return {
    data: data?.tenant || '',
    loading,
    error: null,
  };
};
