import {
  NoteInput,
  CreateOrganizationNoteMutation,
  useCreateOrganizationNoteMutation,
  DataSource,
  GetOrganizationTimelineQuery,
  GetOrganizationTimelineDocument,
  Note,
} from './types';
import { toast } from 'react-toastify';
import client from '../../apollo-client';
import { ApolloCache } from 'apollo-cache';

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
export const useCreateOrganizationNote = ({
  organizationId,
}: Props): Result => {
  const [createOrganizationNoteMutation, { loading, error, data }] =
    useCreateOrganizationNoteMutation();

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { note_CreateForOrganization } }: any,
  ) => {
    const data: GetOrganizationTimelineQuery | null = cache.readQuery({
      query: GetOrganizationTimelineDocument,
      variables: {
        id: organizationId,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationTimelineDocument,
        data: {
          organization: {
            timelineEvents: {
              content: [note_CreateForOrganization],
            },
          },
        },
        variables: { id: organizationId },
      });
      return;
    }

    const newData = {
      organization: [
        ...(data?.organization?.timelineEvents ?? []),
        note_CreateForOrganization,
      ],
    };

    client.writeQuery({
      query: GetOrganizationTimelineDocument,
      data: newData,
      variables: { id: organizationId },
    });
  };
  const handleCreateOrganizationNote: Result['onCreateOrganizationNote'] =
    async (note) => {
      try {
        const response = await createOrganizationNoteMutation({
          variables: { organizationId, input: note },

          optimisticResponse: {
            __typename: 'Mutation',
            note_CreateForOrganization: {
              __typename: 'Note',
              id: 'temp-id',
              appSource: note.appSource || DataSource.Openline,
              html: note.html,
              createdAt: new Date().toISOString(),
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
