import {
  useGetContactPersonalDetailsWithOrganizationsQuery,
  GetContactPersonalDetailsWithOrganizationsQuery,
} from './types';
import { ApolloError } from 'apollo-client';

interface Props {
  id: string;
}

interface Result {
  data:
    | GetContactPersonalDetailsWithOrganizationsQuery['contact']
    | undefined
    | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useContactPersonalDetailsWithOrganizations = ({
  id,
}: Props): Result => {
  const { data, loading, error } =
    useGetContactPersonalDetailsWithOrganizationsQuery({
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
