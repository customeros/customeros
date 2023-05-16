import { ApolloError } from '@apollo/client';
import {GetContactNameByIdQuery, useGetContactNameByIdLazyQuery, useGetContactNameByIdQuery} from './types';

interface Props {
  id: string;
}

interface Result {
  data: GetContactNameByIdQuery['contact'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactNameFromId = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactNameByIdQuery({
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
