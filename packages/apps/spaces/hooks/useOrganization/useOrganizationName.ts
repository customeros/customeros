import { ApolloError } from '@apollo/client';
import { GetOrganizationNameQuery, useGetOrganizationNameQuery } from './types';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationNameQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationName = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationNameQuery({
    variables: { id },
  });

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
    data: data?.organization,
    loading,
    error: null,
  };
};
