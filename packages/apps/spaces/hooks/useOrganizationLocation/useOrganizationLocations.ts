import {
  GetOrganizationLocationsQuery,
  useGetOrganizationLocationsQuery,
} from './types';
import { ApolloError } from '@apollo/client';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationLocationsQuery['organization'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationLocations = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationLocationsQuery({
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
