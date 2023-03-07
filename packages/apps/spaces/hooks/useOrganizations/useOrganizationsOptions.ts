import {
  GetOrganizationsOptionsQuery,
  useGetOrganizationsOptionsQuery,
} from './types';
import { ApolloError } from 'apollo-client';

interface Props {
  id: string;
}

interface Result {
  data:
    | GetOrganizationsOptionsQuery['organizations']['content']
    | undefined
    | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationsOptions = (): Result => {
  const { data, loading, error } = useGetOrganizationsOptionsQuery();

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

  console.log('🏷️ ----- data: ', data);

  return {
    data: data?.organizations.content,
    loading,
    error: null,
  };
};
