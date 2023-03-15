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
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

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
  const [createOrganizationNoteMutation, { loading, error, data }] =
    useCreateOrganizationNoteMutation();
  const { id: userId } = useRecoilValue(userData);

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

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationTimelineDocument,
        data: {
          organization: {
            id: organizationId,
            timelineEvents: [note_CreateForOrganization],
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
          note_CreateForOrganization,
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

          optimisticResponse: {
            __typename: 'Mutation',
            note_CreateForOrganization: {
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
