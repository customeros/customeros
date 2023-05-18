import { ApolloError } from '@apollo/client';
import {
  GetOrganizationNameQuery,
  useGetOrganizationNameLazyQuery,
} from './types';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationNameQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
  onGetOrganizationName: any;
}
export const useOrganizationName = (): Result => {
  const [onGetOrganizationName, { data, loading, error }] =
    useGetOrganizationNameLazyQuery();

  return {
    data: data?.organization,
    loading,
    error: null,
    onGetOrganizationName,
  };
};
