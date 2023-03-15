import {
  NoteInput,
  CreateContactNoteMutation,
  useCreateContactNoteMutation,
  GetOrganizationTimelineQuery,
  GetOrganizationTimelineDocument,
  DataSource,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { DATE_NOW } from '../useOrganizationTimeline/useOrganizationTimeline';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

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
  const [createContactNoteMutation, { loading, error, data }] =
    useCreateContactNoteMutation();
  const { id: userId } = useRecoilValue(userData);

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { note_CreateForOrganization } }: any,
  ) => {
    const data: GetOrganizationTimelineQuery | null = client.readQuery({
      query: GetOrganizationTimelineDocument,
      variables: {
        contactId,
        from: DATE_NOW,
        size: 10,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationTimelineDocument,
        data: {
          organization: {
            contactId,
            timelineEvents: [note_CreateForOrganization],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        timelineEvents: [
          ...(data?.organization?.timelineEvents ?? []),
          note_CreateForOrganization,
        ],
      },
    };

    client.writeQuery({
      query: GetOrganizationTimelineDocument,
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
        optimisticResponse: {
          __typename: 'Mutation',
          note_CreateForContact: {
            __typename: 'Note',
            id: 'temp-id',
            appSource: note.appSource || DataSource.Openline,
            html: note.html,
            createdAt: new Date().toISOString(),
            createdBy: {
              id: userId,
              firstName: '',
              lastName: '',
            },
            updatedAt: '',
            source: DataSource.Openline,
            sourceOfTruth: DataSource.Openline,
          },
        },
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
