import {
  GetOrganizationDetailsQuery,
  useGetOrganizationDetailsQuery,
} from './types';
import { ApolloError } from '@apollo/client';
import { toast } from 'react-toastify';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationDetailsQuery['organization'] | undefined | null;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationDetails = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationDetailsQuery({
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
    toast.error('Something went wrong while loading organization details', {
      toastId: `organization-details-${id}-loading-error`,
    });
    return {
      error,
      loading: false,
      data: null,
    };
  }

  return {
    data: data?.organization,
    loading,
    error: null,
  };
};
