import {
  NoteInput,
  useCreateOrganizationNoteMutation,
  GetOrganizationTimelineQuery,
  GetOrganizationTimelineDocument,
} from './types';
import { toast } from 'react-toastify';
import client from '../../apollo-client';
import { gql } from '@apollo/client';
import { ApolloCache } from '@apollo/client/cache';
import { useSetRecoilState } from 'recoil';
import { organizationNewItemsToEdit } from '../../state/organizationDetails';

interface Props {
  organizationId: string;
}

interface Result {
  saving: boolean;
  onCreateOrganizationNote: (input: NoteInput) => void;
}

const NOW_DATE = new Date().toISOString();
export const useCreateOrganizationNote = ({
  organizationId,
}: Props): Result => {
  const setNoteToEditMode = useSetRecoilState(organizationNewItemsToEdit);

  const [createOrganizationNoteMutation, { loading }] =
    useCreateOrganizationNoteMutation({
      onError: () => {
        toast.error('Something went wrong while adding a note', {
          toastId: `note-add-error-${organizationId}`,
        });
      },
      onCompleted: ({ note_CreateForOrganization }) => {
        setNoteToEditMode((itemsInEditMode) => ({
          timelineEvents: [
            ...itemsInEditMode.timelineEvents,
            { id: note_CreateForOrganization.id },
          ],
        }));
      },
    });

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
      mentioned: [],
      noted: [
        {
          ...organizationData,
          organizationName: organizationData.organizationName || '',
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
  const handleCreateOrganizationNote: Result['onCreateOrganizationNote'] = (
    note,
  ) => {
    return createOrganizationNoteMutation({
      variables: { organizationId, input: note },
      update: handleUpdateCacheAfterAddingNote,
    });
  };

  return {
    saving: loading,
    onCreateOrganizationNote: handleCreateOrganizationNote,
  };
};
