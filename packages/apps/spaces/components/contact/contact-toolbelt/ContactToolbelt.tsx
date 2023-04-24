import React from 'react';

import { useCreateContactNote } from '../../../hooks/useNote';
import { useCreateMeetingFromContact } from '../../../hooks/useMeeting';
import { TimelineToolbelt } from '../../ui-kit/molecules';
import { useCreatePhoneCallInteractionEvent } from '../../../hooks/useContact/useCreatePhoneCallInteractionEvent';
import { useRecoilState } from 'recoil';
import { contactNewItemsToEdit } from '../../../state';

interface ToolbeltProps {
  contactId: string;
}

export const ContactToolbelt: React.FC<ToolbeltProps> = ({ contactId }) => {
  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );

  const { onCreateContactNote, saving } = useCreateContactNote({ contactId });
  const { onCreateMeeting } = useCreateMeetingFromContact({ contactId });
  const { onCreatePhoneCallInteractionEvent } =
    useCreatePhoneCallInteractionEvent({ contactId });

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

  const handleCreatePhoneCallInteractionEvent = (data: any) =>
    onCreatePhoneCallInteractionEvent(data).then((response) => {
      if (response?.id) {
        setItemToEditMode({
          timelineEvents: [
            ...itemsInEditMode.timelineEvents,
            { id: response.id },
          ],
        });
      }
    });

  const handleCreateMeeting = () => onCreateMeeting();

  return (
    <TimelineToolbelt
      onCreateMeeting={handleCreateMeeting}
      onCreateNote={handleCreateNote}
      onLogPhoneCall={handleCreatePhoneCallInteractionEvent}
    />
  );
};
