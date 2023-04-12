import { ApolloError } from 'apollo-client';
import { GetContactTagsQuery, useGetContactTagsQuery } from './types';
import { toast } from 'react-toastify';

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
    toast.error('Something went wrong while loading contact tags', {
      toastId: `get-contact-tags-loading-error`,
    });

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
