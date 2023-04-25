import {
  NoteInput,
  CreateContactNoteMutation,
  useCreateContactNoteMutation,
  GetContactTimelineQuery,
  GetContactTimelineDocument,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';
import { useRecoilState, useSetRecoilState } from 'recoil';
import { contactNewItemsToEdit } from '../../state';

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
  const [createContactNoteMutation, { loading }] =
    useCreateContactNoteMutation();

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { note_CreateForContact } }: any,
  ) => {
    console.log('🏷️ ----- note_CreateForContact: ', note_CreateForContact);
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTimelineDocument,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });
    // @ts-expect-error fix function type
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
    console.log('🏷️ ----- note: ', note);
    try {
      const response = await createContactNoteMutation({
        variables: { contactId, input: note },
        // @ts-expect-error this should not result in error, debug later
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
