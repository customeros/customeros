import { GetContactLocationsQuery, useGetContactLocationsQuery } from './types';
import { ApolloError } from '@apollo/client';

interface Props {
  id: string;
}

interface Result {
  data: GetContactLocationsQuery['contact'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactLocations = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactLocationsQuery({
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
    data: data?.contact ?? null,
    loading,
    error: null,
  };
};
