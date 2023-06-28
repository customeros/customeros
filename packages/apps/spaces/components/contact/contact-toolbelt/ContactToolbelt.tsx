import React from 'react';

import { useCreateContactNote } from '@spaces/hooks/useNote';
import { useCreateMeetingFromContact } from '@spaces/hooks/useMeeting';
import { TimelineToolbelt } from '@spaces/molecules/timeline-toolbelt';
import { useRecoilValue } from 'recoil';
import { userData } from '../../../state';
import { toast } from 'react-toastify';

interface ToolbeltProps {
  contactId: string;
  isSkewed: boolean;
}

export const ContactToolbelt: React.FC<ToolbeltProps> = ({
  contactId,
  isSkewed,
}) => {
  const { id } = useRecoilValue(userData);
  const { onCreateContactNote } = useCreateContactNote({ contactId });
  const { onCreateMeeting } = useCreateMeetingFromContact({ contactId });

  const handleCreateMeeting = () => {
    if (!id) {
      toast.error('Meeting could not be created, please try again later');
      return;
    }
    return onCreateMeeting(id);
  };

  return (
    <TimelineToolbelt
      showPhoneCallButton
      onCreateMeeting={handleCreateMeeting}
      onCreateNote={onCreateContactNote}
      isSkewed={isSkewed}
    />
  );
};
