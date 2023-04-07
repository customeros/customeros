import {
  GetContactCommunicationChannelsQuery,
  useAddEmailToContactMutation,
} from './types';
import {
  EmailInput,
  GetContactCommunicationChannelsDocument,
} from '../../graphQL/__generated__/generated';
import { AddEmailToContactMutation } from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Result {
  onAddEmailToContact: (
    input: EmailInput,
  ) => Promise<AddEmailToContactMutation['emailMergeToContact'] | null>;
}
export const useAddEmailToContactEmail = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [addEmailToContactMutation, { loading, error, data }] =
    useAddEmailToContactMutation();
  const handleUpdateCacheAfterAddingPhoneCall = (
    cache: ApolloCache<any>,
    { data: { emailMergeToContact } }: any,
  ) => {
    const data: GetContactCommunicationChannelsQuery | null = client.readQuery({
      query: GetContactCommunicationChannelsDocument,
      variables: {
        id: contactId,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetContactCommunicationChannelsDocument,
        variables: {
          id: contactId,
        },
        data: {
          contact: {
            id: contactId,
            emails: [emailMergeToContact],
          },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        emails: [
          ...(data.contact?.emails || []),
          { ...emailMergeToContact, email: emailMergeToContact.email || '' },
        ],
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

  const handleAddEmailToContact: Result['onAddEmailToContact'] = async (
    email,
  ) => {
    try {
      console.log('🏷️ ----- email: ', email);
      const optimisticItem = { id: 'optimistic-id', ...email };
      console.log('🏷️ ----- email: ', email);
      const response = await addEmailToContactMutation({
        variables: { contactId, input: email },
        optimisticResponse: {
          emailMergeToContact: {
            __typename: 'Email',
            ...optimisticItem,
            primary: optimisticItem?.primary || false,
          },
        },
        // @ts-expect-error fixme
        update: handleUpdateCacheAfterAddingPhoneCall,
      });
      return response.data?.emailMergeToContact ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding email');

      return null;
    }
  };

  return {
    onAddEmailToContact: handleAddEmailToContact,
  };
};
