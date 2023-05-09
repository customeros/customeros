import {
  NoteInput,
  CreateOrganizationNoteMutation,
  useCreateOrganizationNoteMutation,
  GetOrganizationTimelineQuery,
  GetOrganizationTimelineDocument,
} from './types';
import { toast } from 'react-toastify';
import client from '../../apollo-client';
import { gql } from '@apollo/client';
import { ApolloCache } from '@apollo/client/cache';

interface Props {
  organizationId: string;
}

interface Result {
  saving: boolean;
  onCreateOrganizationNote: (
    input: NoteInput,
  ) => Promise<
    CreateOrganizationNoteMutation['note_CreateForOrganization'] | null
  >;
}

const NOW_DATE = new Date().toISOString();
export const useCreateOrganizationNote = ({
  organizationId,
}: Props): Result => {
  const [createOrganizationNoteMutation, { loading }] =
    useCreateOrganizationNoteMutation();

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { note_CreateForOrganization } }: any,
  ) => {
    const data: GetOrganizationTimelineQuery | null = client.readQuery({
      query: GetOrganizationTimelineDocument,
      variables: {
        organizationId,
        from: NOW_DATE,
        size: 10,
      },
    });

    const normalizedId = cache.identify({
      id: organizationId,
      __typename: 'Organization',
    });
    const organizationData = client.readFragment({
      id: normalizedId,
      fragment: gql`
        fragment organizationName on Organization {
          id
          name
        }
      `,
    });
    const newNoteWithNoted = {
      ...note_CreateForOrganization,
      noted: [
        {
          ...organizationData,
        },
      ],
    };
    if (data === null) {
      client.writeQuery({
        query: GetOrganizationTimelineDocument,
        data: {
          organization: {
            id: organizationId,
            timelineEvents: [newNoteWithNoted],
          },
          variables: { organizationId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        timelineEvents: [
          ...(data?.organization?.timelineEvents ?? []),
          newNoteWithNoted,
        ],
      },
    };

    client.writeQuery({
      query: GetOrganizationTimelineDocument,
      data: newData,
      variables: {
        organizationId,
        from: NOW_DATE,
        size: 10,
      },
    });
  };
  const handleCreateOrganizationNote: Result['onCreateOrganizationNote'] =
    async (note) => {
      try {
        const response = await createOrganizationNoteMutation({
          variables: { organizationId, input: note },
          update: handleUpdateCacheAfterAddingNote,
        });
        if (response.data) {
          toast.success('Note added!', {
            toastId: `note-added-${response.data?.note_CreateForOrganization.id}`,
          });
        }
        return response.data?.note_CreateForOrganization ?? null;
      } catch (err) {
        toast.error('Something went wrong while adding a note', {
          toastId: `note-add-error-${organizationId}`,
        });
        return null;
      }
    };

  return {
    saving: loading,
    onCreateOrganizationNote: handleCreateOrganizationNote,
  };
};
