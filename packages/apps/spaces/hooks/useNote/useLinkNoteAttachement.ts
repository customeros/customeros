import {
  GetContactTimelineDocument,
  NOW_DATE,
  useNoteLinkAttachmentMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import {
  GetContactTimelineQuery,
  NoteLinkAttachmentMutation,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  noteId: string;
  contactId?: string;
}

export interface Result {
  onLinkNoteAttachment: (
    attachmentId: string,
  ) => Promise<NoteLinkAttachmentMutation['note_LinkAttachment'] | null>;
}

export const useLinkNoteAttachment = ({ noteId, contactId }: Props): Result => {
  const [linkNoteAttachmentMutation, { loading, error, data }] =
    useNoteLinkAttachmentMutation();
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

  const handleLinkNoteAttachment: Result['onLinkNoteAttachment'] = async (
    attachmentId,
  ) => {
    try {
      const response = await linkNoteAttachmentMutation({
        variables: {
          noteId,
          attachmentId,
        },
        //update: handleUpdateCacheAfterAddingNote,
      });

      toast.success(`Added attachment to note`);
      return response.data?.note_LinkAttachment ?? null;
    } catch (err) {
      console.error(err);
      toast.error(`Something went wrong while attaching file to the note`);
      return null;
    }
  };

  return {
    onLinkNoteAttachment: handleLinkNoteAttachment,
  };
};
