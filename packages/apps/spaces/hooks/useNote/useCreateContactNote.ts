import {
  NoteInput,
  CreateContactNoteMutation,
  useCreateContactNoteMutation,
  GetContactTimelineQuery,
  GetContactTimelineDocument,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';

interface Props {
  contactId: string;
}

interface Result {
  saving: boolean;
  onCreateContactNote: (
    input: NoteInput,
  ) => Promise<CreateContactNoteMutation['note_CreateForContact'] | null>;
}

const NOW_DATE = new Date().toISOString();

export const useCreateContactNote = ({ contactId }: Props): Result => {
  const [createContactNoteMutation, { loading }] = useCreateContactNoteMutation(
    { fetchPolicy: 'no-cache' },
  );

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { note_CreateForContact } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTimelineDocument,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });

    const normalizedId = cache.identify({
      id: contactId,
      __typename: 'Contact',
    });
    const contactData = client.readFragment({
      id: normalizedId,
      fragment: gql`
        fragment ContactName on Contact {
          id
          name
          firstName
          lastName
        }
      `,
    });
    const newNoteWithNoted = {
      ...note_CreateForContact,
      noted: [
        {
          ...contactData,
        },
      ],
    };

    client.writeFragment({
      id: `Note:${note_CreateForContact.id}`,
      fragment: gql`
        fragment NoteF on Note {
          id
          html
          createdAt
          source
          noted {
            ... on Organization {
              id
              organizationName: name
            }
            ... on Contact {
              firstName
              lastName
            }
          }
          createdBy {
            id
            firstName
            lastName
          }
          includes {
            id
            name
            mimeType
            extension
            size
          }
        }
      `,

      data: {
        ...note_CreateForContact,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [newNoteWithNoted],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [newNoteWithNoted],
      },
    };

    client.writeQuery({
      query: GetContactTimelineDocument,
      data: newData,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });
  };

  const handleCreateContactNote: Result['onCreateContactNote'] = async (
    note,
  ) => {
    try {
      const response = await createContactNoteMutation({
        variables: { contactId, input: note },
        update: handleUpdateCacheAfterAddingNote,
      });
      if (response.data) {
        toast.success('Note added!', {
          toastId: `note-added-${response.data?.note_CreateForContact.id}`,
        });
      }
      return response.data?.note_CreateForContact ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding a note', {
        toastId: `note-add-error-${contactId}`,
      });
      return null;
    }
  };

  return {
    saving: loading,
    onCreateContactNote: handleCreateContactNote,
  };
};
