import React from 'react';
import { useCreateOrganizationNote } from '@spaces/hooks/useNote';
import { useCreateMeetingFromOrganization } from '@spaces/hooks/useMeeting';
import { TimelineToolbelt } from '@spaces/molecules/timeline-toolbelt';
import { useRecoilValue } from 'recoil';
import { userData } from '../../../state';
import { toast } from 'react-toastify';

interface ToolbeltProps {
  organizationId: string;
}

export const OrginizationToolbelt: React.FC<ToolbeltProps> = ({
  organizationId,
}) => {
  const { id } = useRecoilValue(userData);

  const { onCreateOrganizationNote } = useCreateOrganizationNote({
    organizationId,
  });
  const { onCreateMeeting } = useCreateMeetingFromOrganization({
    organizationId,
  });

  const handleCreateMeeting = () => {
    if (!id) {
      toast.error('Meeting could not be created, please try again later');
      return;
    }
    onCreateMeeting(id);
  };
  return (
    <TimelineToolbelt
      onCreateMeeting={handleCreateMeeting}
      onCreateNote={onCreateOrganizationNote}
      isSkewed
    />
  );
};
