import {
  GetOrganizationCustomFieldsQuery,
  useGetOrganizationCustomFieldsQuery,
} from './types';
import { ApolloError } from '@apollo/client';
import { toast } from 'react-toastify';

interface Props {
  id: string;
}

interface Result {
  data: GetOrganizationCustomFieldsQuery['organization'] | null | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationCustomFields = ({ id }: Props): Result => {
  const { data, loading, error } = useGetOrganizationCustomFieldsQuery({
    variables: { id },
    nextFetchPolicy: 'cache-first',
  });

  if (loading) {
    return {
      loading: true,
      error: null,
      data: null,
    };
  }

  if (error) {
    toast.error('Something went wrong while loading custom fields', {
      toastId: `organization-custom-fields-${id}-loading-error`,
    });
    return {
      error,
      loading: false,
      data: null,
    };
  }

  return {
    data: data?.organization
      ? {
          ...data?.organization,
          // emails: [...(data?.organization?.emails ?? [])]?.sort((a, b) =>
          //   a.primary === b.primary ? 0 : a.primary ? -1 : 1,
          // ),
          // phoneNumbers: [...(data?.organization?.phoneNumbers ?? [])]?.sort(
          //   //@ts-expect-error fixme
          //   (a, b) => (a.primary === b.primary ? 0 : a.primary ? -1 : 1),
          // ),
        }
      : null,
    loading,
    error: null,
  };
};
