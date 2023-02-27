import {
  GetOrganizationDetailsQuery,
  useGetOrganizationDetailsQuery,
} from './types';
import { ApolloError } from 'apollo-client';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationDetailsQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationDetails = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationDetailsQuery({
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
