import {
  Email,
  EmailUpdateInput,
  GetContactCommunicationChannelsDocument,
  GetContactCommunicationChannelsQuery,
} from '../../graphQL/__generated__/generated';
import {
  UpdateContactEmailMutation,
  useUpdateContactEmailMutation,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Result {
  onUpdateContactEmail: (
    input: EmailUpdateInput,
  ) => Promise<UpdateContactEmailMutation['emailUpdateInContact'] | null>;
}
export const useUpdateContactEmail = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [updateContactNoteMutation, { loading, error, data }] =
    useUpdateContactEmailMutation();
  const handleUpdateCacheAfterMutation = (
    cache: ApolloCache<any>,
    { data: { emailUpdateInContact } }: any,
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
            emails: [emailUpdateInContact],
          },
        },
      });
      return;
    }

    const newEmailData = (data.contact?.emails || []).map((oldEmail) =>
      oldEmail.id === emailUpdateInContact.id ? emailUpdateInContact : oldEmail,
    );
    const newData = {
      contact: {
        ...data.contact,
        emails: newEmailData,
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

  const handleUpdateContactEmail: Result['onUpdateContactEmail'] = async (
    input,
  ) => {
    try {
      const response = await updateContactNoteMutation({
        variables: { input, contactId },
        optimisticResponse: {
          emailUpdateInContact: {
            __typename: 'Email',
            ...input,
            primary: input.primary || false,
          },
        },
        // @ts-expect-error fixme
        update: handleUpdateCacheAfterMutation,
      });

      return response.data?.emailUpdateInContact ?? null;
    } catch (err) {
      toast.error('Something went wrong while updating contact email');
      return null;
    }
  };

  return {
    onUpdateContactEmail: handleUpdateContactEmail,
  };
};
