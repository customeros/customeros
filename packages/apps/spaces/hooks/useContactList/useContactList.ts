import { ApolloError } from '@apollo/client';
import { useGetContactListLazyQuery } from './types';
import { GetContactListQueryVariables } from '../../graphQL/__generated__/generated';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onLoadContactList: ({
    variables,
  }: {
    variables: GetContactListQueryVariables;
  }) => Promise<any>;
}
export const useContactList = (): Result => {
  const [onLoadContactList, { data, loading, error }] =
    useGetContactListLazyQuery();

  return {
    onLoadContactList,
    loading,
    error,
  };
};
