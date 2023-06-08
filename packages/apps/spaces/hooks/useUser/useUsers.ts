import { ApolloError, QueryResult } from '@apollo/client';
import { Filter, Pagination } from './types';
import {
  useGetUsersLazyQuery,
  GetUsersQueryVariables,
  GetUsersQuery,
} from '@spaces/graphql';
import { useRecoilState } from 'recoil';
import { ownerListData } from '../../state/userData';
import { useEffect } from 'react';

interface Props {
  pagination: Pagination;
  where?: Filter;
}

interface Result {
  loading: boolean;

  error: ApolloError | null;
  onLoadUsers: ({
    variables,
  }: {
    variables: GetUsersQueryVariables;
  }) => Promise<QueryResult<GetUsersQuery, GetUsersQueryVariables>>;
}
export const useUsers = (): Result => {
  const [loadUsers, { loading, error }] = useGetUsersLazyQuery();
  const [ownerListResult, setOwnersList] = useRecoilState(ownerListData);

  useEffect(() => {
    if (!ownerListResult.ownerList.length) {
      loadUsers({
        variables: {
          pagination: { page: 0, limit: 100 },
        },
      }).then((res) => {
        const ownerList = (res.data?.users?.content ?? []).map((data) => ({
          label: `${data?.firstName} ${data?.lastName}`,
          value: data.id,
        }));
        setOwnersList({ ownerList });
      });
    }
  }, []);

  return {
    loading,
    onLoadUsers: loadUsers,
    error: null,
  };
};
