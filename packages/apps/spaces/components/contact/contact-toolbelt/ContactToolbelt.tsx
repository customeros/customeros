import React from 'react';

import { useCreateContactNote } from '../../../hooks/useNote';
import { useCreateMeetingFromContact } from '../../../hooks/useMeeting';
import { TimelineToolbelt } from '../../ui-kit/molecules';
import { useCreatePhoneCallInteractionEvent } from '../../../hooks/useContact/useCreatePhoneCallInteractionEvent';

interface ToolbeltProps {
  contactId: string;
}

export const ContactToolbelt: React.FC<ToolbeltProps> = ({ contactId }) => {
  const { onCreateContactNote, saving } = useCreateContactNote({ contactId });
  const { onCreateMeeting } = useCreateMeetingFromContact({ contactId });
  const { onCreatePhoneCallInteractionEvent } =
    useCreatePhoneCallInteractionEvent({ contactId });

  return (
    <TimelineToolbelt
      onCreateMeeting={onCreateMeeting}
      onCreateNote={onCreateContactNote}
      onLogPhoneCall={onCreatePhoneCallInteractionEvent}
    />
  );
};
