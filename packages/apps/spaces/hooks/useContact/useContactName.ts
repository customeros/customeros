import { ApolloError } from 'apollo-client';
import { GetContactNameQuery, useGetContactNameQuery } from './types';

interface Props {
  id: string;
}

interface Result {
  data: GetContactNameQuery['contact'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactName = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactNameQuery({
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
    data: data?.contact,
    loading,
    error: null,
  };
};
