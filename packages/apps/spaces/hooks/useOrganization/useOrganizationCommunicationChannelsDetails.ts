import {
  GetOrganizationCommunicationChannelsQuery,
  useGetOrganizationCommunicationChannelsQuery,
} from './types';
import { ApolloError } from '@apollo/client';
import { toast } from 'react-toastify';

interface Props {
  id: string;
}

interface Result {
  data:
    | GetOrganizationCommunicationChannelsQuery['organization']
    | null
    | undefined;
  loading: boolean;
  error: ApolloError | null;
}
export const useOrganizationCommunicationChannelsDetails = ({
  id,
}: Props): Result => {
  const { data, loading, error } = useGetOrganizationCommunicationChannelsQuery(
    {
      variables: { id },
      nextFetchPolicy: 'cache-first',
    },
  );

  if (loading) {
    return {
      loading: true,
      error: null,
      data: null,
    };
  }

  if (error) {
    toast.error('Something went wrong while loading communication channels', {
      toastId: `organization-communication-channels-${id}-loading-error`,
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
          emails: [...(data?.organization?.emails ?? [])]?.sort((a, b) =>
            a.primary === b.primary ? 0 : a.primary ? -1 : 1,
          ),
          phoneNumbers: [...(data?.organization?.phoneNumbers ?? [])]?.sort(
            //@ts-expect-error fixme
            (a, b) => (a.primary === b.primary ? 0 : a.primary ? -1 : 1),
          ),
        }
      : null,
    loading,
    error: null,
  };
};
