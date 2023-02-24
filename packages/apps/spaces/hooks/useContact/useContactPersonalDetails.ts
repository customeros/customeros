import {
  GetContactPersonalDetailsQuery,
  useGetContactPersonalDetailsQuery,
} from './types';
import { ApolloError } from 'apollo-client';

interface Props {
  id: string;
}

interface Result {
  data: GetContactPersonalDetailsQuery['contact'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactPersonalDetails = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactPersonalDetailsQuery({
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
