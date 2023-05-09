import {
  GetContactPersonalDetailsQuery,
  useGetContactPersonalDetailsQuery,
} from './types';
import { ApolloError } from '@apollo/client';
import { toast } from 'react-toastify';

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
