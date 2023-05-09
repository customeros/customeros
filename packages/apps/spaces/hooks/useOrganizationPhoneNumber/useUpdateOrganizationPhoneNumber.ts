import {
  PhoneNumberUpdateInput,
  UpdateOrganizationPhoneNumberMutation,
  GetOrganizationCommunicationChannelsQuery,
  useUpdateOrganizationPhoneNumberMutation,
  GetOrganizationCommunicationChannelsDocument,
} from './types';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Result {
  onUpdateOrganizationPhoneNumber: (
    input: PhoneNumberUpdateInput,
  ) => Promise<
    | UpdateOrganizationPhoneNumberMutation['phoneNumberUpdateInOrganization']
    | null
  >;
}
export const useUpdateOrganizationPhoneNumber = ({
  organizationId,
}: {
  organizationId: string;
}): Result => {
  const [updateOrganizationNoteMutation, { loading, error, data }] =
    useUpdateOrganizationPhoneNumberMutation();
  const handleUpdateCacheAfterUpdatingPhoneNumber = (
    cache: ApolloCache<any>,
    { data: { phoneNumberUpdateInOrganization } }: any,
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
            phoneNumbers: [phoneNumberUpdateInOrganization],
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        phoneNumbers: (data.organization?.phoneNumbers || []).map((e) =>
          e.id === phoneNumberUpdateInOrganization.id
            ? { ...e, ...phoneNumberUpdateInOrganization }
            : {
                ...e,
                primary: phoneNumberUpdateInOrganization.primary
                  ? false
                  : //@ts-expect-error revisit later
                    !!e?.primary,
              },
        ),
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
  const handleUpdateOrganizationPhoneNumber: Result['onUpdateOrganizationPhoneNumber'] =
    async (input) => {
      const payload = {
        ...input,
      };
      try {
        const response = await updateOrganizationNoteMutation({
          variables: { input: payload, organizationId },
          update: handleUpdateCacheAfterUpdatingPhoneNumber,
        });

        return response.data?.phoneNumberUpdateInOrganization ?? null;
      } catch (err) {
        console.error(err);
        toast.error(
          'Something went wrong while updating email! Please contact us or try again later',
          {
            toastId: `update-organization-email-error-${input.id}-${organizationId}`,
          },
        );

        return null;
      }
    };

  return {
    onUpdateOrganizationPhoneNumber: handleUpdateOrganizationPhoneNumber,
  };
};
