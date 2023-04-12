import {
  GetContactCommunicationChannelsDocument,
  GetContactCommunicationChannelsQuery,
  PhoneNumber,
  PhoneNumberUpdateInput,
} from '../../graphQL/__generated__/generated';
import {
  UpdateContactPhoneNumberMutation,
  useUpdateContactPhoneNumberMutation,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Result {
  onUpdateContactPhoneNumber: (
    input: PhoneNumberUpdateInput,
  ) => Promise<
    UpdateContactPhoneNumberMutation['phoneNumberUpdateInContact'] | null
  >;
}
export const useUpdateContactPhoneNumber = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [updateContactNoteMutation, { loading, error, data }] =
    useUpdateContactPhoneNumberMutation();
  const handleUpdateCacheAfterAddingPhoneNumber = (
    cache: ApolloCache<any>,
    { data: { phoneNumberUpdateInContact } }: any,
  ) => {
    const data: GetContactCommunicationChannelsQuery | null = client.readQuery({
      query: GetContactCommunicationChannelsDocument,
      variables: {
        id: contactId,
      },
    });
    console.log('🏷️ ----- data: ', data);
    if (data === null) {
      client.writeQuery({
        query: GetContactCommunicationChannelsDocument,
        variables: {
          id: contactId,
        },
        data: {
          contact: {
            id: contactId,
            phoneNumbers: [phoneNumberUpdateInContact],
          },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        phoneNumbers: (data.contact?.phoneNumbers || []).map((e) =>
          e.id === phoneNumberUpdateInContact.id
            ? { ...e, ...phoneNumberUpdateInContact }
            : {
                ...e,
                primary: phoneNumberUpdateInContact.primary ? false : e.primary,
              },
        ),
      },
    };

    client.writeQuery({
      query: GetContactCommunicationChannelsDocument,
      data: newData,
      variables: {
        id: contactId,
      },
    });
  };
  const handleUpdateContactPhoneNumber: Result['onUpdateContactPhoneNumber'] =
    async (input) => {
      const payload = {
        ...input,
      };
      try {
        const response = await updateContactNoteMutation({
          variables: { input: payload, contactId },
          // @ts-expect-error fixme
          update: handleUpdateCacheAfterAddingPhoneNumber,
        });

        return response.data?.phoneNumberUpdateInContact ?? null;
      } catch (err) {
        console.error(err);
        toast.error(
          'Something went wrong while updating phone number! Please contact us or try again later',
          {
            toastId: `update-contact-phone-error-${input.id}-${contactId}`,
          },
        );
        return null;
      }
    };

  return {
    onUpdateContactPhoneNumber: handleUpdateContactPhoneNumber,
  };
};
