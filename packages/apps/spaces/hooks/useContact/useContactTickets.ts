import { ApolloError } from 'apollo-client';
import { GetContactTagsQuery, useGetContactTagsQuery } from './types';
import {
  GetContactTicketsQuery,
  useGetContactTicketsQuery,
} from '../../graphQL/__generated__/generated';

interface Props {
  id: string;
}

interface Result {
  data: GetContactTicketsQuery['contact'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactTickets = ({ id }: Props): Result => {
  const { data, loading, error } = useGetContactTicketsQuery({
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
