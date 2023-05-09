import { ApolloError } from '@apollo/client';
import {
  OrganizationContactsFragment,
  useGetOrganizationContactsQuery,
} from './types';

interface Props {
  id: string;
}

interface Result {
  data: OrganizationContactsFragment['contacts']['content'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationContacts = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationContactsQuery({
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
    data: data?.organization?.contacts?.content ?? [],
    loading,
    error: null,
  };
};
