import { ApolloError } from 'apollo-client';
import { GetContactTagsQuery, useGetContactTagsQuery } from './types';

interface Props {
  id: string;
}

interface Result {
  data: GetContactTagsQuery['contact'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactTags = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactTagsQuery({
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
