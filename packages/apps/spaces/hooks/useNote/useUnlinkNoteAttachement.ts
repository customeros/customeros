import {
  GetContactTimelineDocument,
  NOW_DATE,
  useNoteUnlinkAttachmentMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import {
  GetContactTimelineQuery,
  NoteUnlinkAttachmentMutation,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  noteId: string;
  contactId?: string;
}

export interface Result {
  onUnlinkNoteAttachment: (
    fileId: string,
  ) => Promise<NoteUnlinkAttachmentMutation['note_UnlinkAttachment'] | null>;
}

export const useUnlinkNoteAttachment = ({
  noteId,
  contactId,
}: Props): Result => {
  const [unlinkNoteAttachmentMutation, { loading, error, data }] =
    useNoteUnlinkAttachmentMutation();
  const loggedInUserData = useRecoilValue(userData);

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { note_Create } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTimelineDocument,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [note_Create],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [...(data.contact?.timelineEvents || []), note_Create],
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

  const handleUnlinkNoteAttachment: Result['onUnlinkNoteAttachment'] = async (
    attachmentId,
  ) => {
    try {
      const response = await unlinkNoteAttachmentMutation({
        variables: {
          noteId,
          attachmentId,
        },

        //update: handleUpdateCacheAfterAddingNote,
      });

      return response.data?.note_UnlinkAttachment ?? null;
    } catch (err) {
      console.error(err);
      toast.error(
        `Something went wrong while adding draft note to the timeline`,
      );
      return null;
    }
  };

  return {
    onUnlinkNoteAttachment: handleUnlinkNoteAttachment,
  };
};
