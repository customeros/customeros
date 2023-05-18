import { ApolloError } from '@apollo/client';
import { useGetContactNameByIdLazyQuery } from './types';

interface Result {
  loading: boolean;
  error: ApolloError | undefined;
  onGetContactNameById: any;
}
export const useContactNameFromId = (): Result => {
  const [onGetContactNameById, { loading, error }] =
    useGetContactNameByIdLazyQuery();

  return {
    loading,
    error,
    onGetContactNameById,
  };
};
