import { ApolloError } from '@apollo/client';
import { GetOrganizationSubsidiariesQuery, useGetOrganizationSubsidiariesQuery } from './types';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationSubsidiariesQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationSubsidiaries = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationSubsidiariesQuery({
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
