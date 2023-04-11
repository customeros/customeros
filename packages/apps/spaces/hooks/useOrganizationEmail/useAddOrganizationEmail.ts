import {
  useAddEmailToOrganizationMutation,
  AddEmailToOrganizationMutation,
} from './types';
import {
  EmailInput,
  GetOrganizationCommunicationChannelsDocument,
  GetOrganizationCommunicationChannelsQuery,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Result {
  onAddEmailToOrganization: (
    input: EmailInput,
  ) => Promise<
    AddEmailToOrganizationMutation['emailMergeToOrganization'] | null
  >;
}
export const useAddEmailToOrganizationEmail = ({
  organizationId,
}: {
  organizationId: string;
}): Result => {
  const [addEmailToOrganizationMutation, { loading, error, data }] =
    useAddEmailToOrganizationMutation();
  const handleUpdateCacheAfterAddingPhoneCall = (
    cache: ApolloCache<any>,
    { data: { emailMergeToOrganization } }: any,
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
            emails: [emailMergeToOrganization],
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        emails: [
          ...(data.organization?.emails || []),
          {
            ...emailMergeToOrganization,
            email: emailMergeToOrganization.email || '',
          },
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

  const handleAddEmailToOrganization: Result['onAddEmailToOrganization'] =
    async (email) => {
      try {
        const optimisticItem = { id: 'optimistic-id', ...email };
        const response = await addEmailToOrganizationMutation({
          variables: { organizationId, input: email },
          optimisticResponse: {
            emailMergeToOrganization: {
              __typename: 'Email',
              ...optimisticItem,
              primary: optimisticItem?.primary || false,
            },
          },
          // @ts-expect-error fixme
          update: handleUpdateCacheAfterAddingPhoneCall,
        });
        return response.data?.emailMergeToOrganization ?? null;
      } catch (err) {
        toast.error('Something went wrong while adding email');

        return null;
      }
    };

  return {
    onAddEmailToOrganization: handleAddEmailToOrganization,
  };
};
