import {
  AddPhoneToOrganizationMutation,
  GetOrganizationCommunicationChannelsQuery,
  useAddPhoneToOrganizationMutation,
  GetOrganizationCommunicationChannelsDocument,
} from './types';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Props {
  organizationId: string;
}

interface Result {
  onCreateOrganizationPhoneNumber: (
    input: any, //FIXME
  ) => Promise<
    AddPhoneToOrganizationMutation['phoneNumberMergeToOrganization'] | null
  >;
  loading: boolean;
}
export const useCreateOrganizationPhoneNumber = ({
  organizationId,
}: Props): Result => {
  const [createOrganizationPhoneNumberMutation, { loading, error, data }] =
    useAddPhoneToOrganizationMutation();

  const handleUpdateCacheAfterAddingPhoneNumber = (
    cache: ApolloCache<any>,
    { data: { phoneNumberMergeToOrganization } }: any,
  ) => {
    const data: GetOrganizationCommunicationChannelsQuery | null =
      client.readQuery({
        query: GetOrganizationCommunicationChannelsDocument,
        variables: {
          id: organizationId,
        },
      });
    if (data === null) {
      client.writeQuery({
        query: GetOrganizationCommunicationChannelsDocument,
        variables: {
          id: organizationId,
        },
        data: {
          organization: {
            id: organizationId,
            phoneNumbers: [phoneNumberMergeToOrganization],
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        phoneNumbers: [
          ...(data.organization?.phoneNumbers || []),
          phoneNumberMergeToOrganization,
        ],
      },
    };
    client.writeQuery({
      query: GetOrganizationCommunicationChannelsDocument,
      data: newData,
      variables: {
        id: organizationId,
      },
    });
  };

  const handleCreateOrganizationPhoneNumber: Result['onCreateOrganizationPhoneNumber'] =
    async (input) => {
      try {
        const response = await createOrganizationPhoneNumberMutation({
          variables: { organizationId, input },
          refetchQueries: ['GetOrganizationCommunicationChannels'],
          awaitRefetchQueries: true,
          optimisticResponse: {
            phoneNumberMergeToOrganization: {
              __typename: 'PhoneNumber',
              ...input,
              id: 'optimistic-id',
              primary: input?.primary || false,
            },
          },
          // @ts-expect-error fixme
          update: handleUpdateCacheAfterAddingPhoneNumber,
        });
        return response.data?.phoneNumberMergeToOrganization ?? null;
      } catch (err) {
        toast.error('Something went wrong while adding phone number');
        console.error(err);
        return null;
      }
    };

  return {
    onCreateOrganizationPhoneNumber: handleCreateOrganizationPhoneNumber,
    loading,
  };
};
