import { ApolloError } from 'apollo-client';
import {
  GetActionsForContactQuery,
  GetContactConversationsQuery,
  Pagination,
  useGetActionsForContactQuery,
  useGetContactConversationsQuery,
} from '../../graphQL/__generated__/generated';

interface Props {
  id: string;
}

interface Result {
  data: GetActionsForContactQuery['contact'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}

export const useContactActions = ({ id }: Props): Result => {
  const from = new Date(1970, 0, 1).toISOString();
  const to = new Date().toISOString();
  const { data, loading, error } = useGetActionsForContactQuery({
    variables: { id, from, to },
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
  console.log('data loaded for actions');
  return {
    data: data?.contact ?? null,
    loading,
    error: null,
  };
};
