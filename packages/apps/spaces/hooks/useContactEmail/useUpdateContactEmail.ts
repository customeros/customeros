import {
  GetContactCommunicationChannelsQuery,
  UpdateContactEmailMutation,
  useUpdateContactEmailMutation,
} from '../useContact/types';
import {
  EmailUpdateInput,
  GetContactCommunicationChannelsDocument,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

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

    const newData = {
      contact: {
        ...data.contact,
        emails: [...(data.contact?.emails || []), { ...emailUpdateInContact }],
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

  const handleUpdateCacheAfterAddingEmail = (
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

    const newData = {
      contact: {
        ...data.contact,
        emails: (data.contact?.emails || []).map((e) =>
          e.id === emailUpdateInContact.id
            ? { ...e, ...emailUpdateInContact }
            : {
                ...e,
                primary: emailUpdateInContact.primary ? false : e.primary,
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
  const handleUpdateContactEmail: Result['onUpdateContactEmail'] = async (
    input,
  ) => {
    try {
      const response = await updateContactEmailMutation({
        variables: { input: { ...input }, contactId },
        //@ts-expect-error fixme
        update: handleUpdateCacheAfterAddingEmail,
      });

      return response.data?.emailUpdateInContact ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong while updating email! Please contact us or try again later',
        {
          toastId: `update-contact-email-error-${input.id}-${contactId}`,
        },
      );
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateContactEmail: handleUpdateContactEmail,
  };
};
