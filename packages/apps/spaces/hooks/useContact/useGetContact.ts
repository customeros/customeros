import { ApolloError } from '@apollo/client';
import { toast } from 'react-toastify';
import { GetContactQuery, useGetContactQuery } from '@spaces/graphql';

interface Props {
  id: string;
}

interface Result {
  data: GetContactQuery['contact'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContact = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactQuery({
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
    toast.error('Something went wrong while loading contact personal details', {
      toastId: `get-contact-personal-details-query-error`,
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
