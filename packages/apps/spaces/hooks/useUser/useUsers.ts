import { ApolloError } from '@apollo/client';
import { Filter, GetUsersQuery, Pagination, useGetUsersQuery } from './types';
import {
  GetContactMentionSuggestionsQueryVariables,
  GetUsersLazyQueryHookResult,
  GetUsersQueryVariables,
  useGetUsersLazyQuery,
} from '../../graphQL/__generated__/generated';

interface Props {
  pagination: Pagination;
  where?: Filter;
}

interface Result {
  data: GetUsersQuery['users'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
  onLoadUsers: ({
    variables,
  }: {
    variables: GetUsersQueryVariables;
  }) => Promise<any>;
}
export const useUsers = (): Result => {
  const [loadUsers, { data, loading, error }] = useGetUsersLazyQuery();

  return {
    data: data?.users,
    loading,
    error: null,
    onLoadUsers: loadUsers,
  };
};
