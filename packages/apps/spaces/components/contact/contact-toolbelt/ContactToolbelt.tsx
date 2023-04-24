import React from 'react';

import { useCreateContactNote } from '../../../hooks/useNote';
import { useCreateMeetingFromContact } from '../../../hooks/useMeeting';
import { TimelineToolbelt } from '../../ui-kit/molecules';
import { useCreatePhoneCallInteractionEvent } from '../../../hooks/useContact/useCreatePhoneCallInteractionEvent';
import { useRecoilState, useRecoilValue } from 'recoil';
import { contactNewItemsToEdit, userData } from '../../../state';
import { useUser } from '../../../hooks/useUser';
import { toast } from 'react-toastify';

interface ToolbeltProps {
  contactId: string;
}

export const ContactToolbelt: React.FC<ToolbeltProps> = ({ contactId }) => {
  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );
  const { identity: userEmail } = useRecoilValue(userData);
  const { data, loading, error } = useUser({ email: userEmail });
  const { onCreateContactNote, saving } = useCreateContactNote({ contactId });
  const { onCreateMeeting } = useCreateMeetingFromContact({ contactId });

  const handleCreateNote = (data: any) =>
    onCreateContactNote(data).then((response) => {
      if (response?.id) {
        setItemToEditMode({
          timelineEvents: [
            ...itemsInEditMode.timelineEvents,
            { id: response.id },
          ],
        });
      }
    });

  const handleCreateMeeting = () => {
    if (!data?.id) {
      toast.error('Meeting could not be created, please try again later');
      return;
    }
    return onCreateMeeting(data?.id);
  };

  return (
    <TimelineToolbelt
      onCreateMeeting={handleCreateMeeting}
      onCreateNote={handleCreateNote}
    />
  );
};
