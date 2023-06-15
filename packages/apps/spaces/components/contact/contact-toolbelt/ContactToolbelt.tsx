import React from 'react';

import { useCreateContactNote } from '@spaces/hooks/useNote';
import { useCreateMeetingFromContact } from '@spaces/hooks/useMeeting';
import { TimelineToolbelt } from '@spaces/molecules/timeline-toolbelt';
import { useRecoilValue } from 'recoil';
import { userData } from '../../../state';
import { useUser } from '@spaces/hooks/useUser';
import { toast } from 'react-toastify';

interface ToolbeltProps {
  contactId: string;
  isSkewed: boolean;
}

export const ContactToolbelt: React.FC<ToolbeltProps> = ({
  contactId,
  isSkewed,
}) => {
  const { identity: userEmail } = useRecoilValue(userData);
  const { data } = useUser({ email: userEmail });
  const { onCreateContactNote } = useCreateContactNote({ contactId });
  const { onCreateMeeting } = useCreateMeetingFromContact({ contactId });

  const handleCreateMeeting = () => {
    if (!data?.id) {
      toast.error('Meeting could not be created, please try again later');
      return;
    }
    return onCreateMeeting(data?.id);
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
