import {
  GetOrganizationOwnerQuery,
  useGetOrganizationOwnerQuery,
} from './types';
import { ApolloError } from '@apollo/client';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationOwnerQuery['organization'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationOwner = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationOwnerQuery({
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
    data: data?.organization ?? null,
    loading,
    error: null,
  };
};
