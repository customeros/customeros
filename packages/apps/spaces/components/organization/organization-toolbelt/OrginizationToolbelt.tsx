import React from 'react';

import { useCreateOrganizationNote } from '../../../hooks/useNote';
import { useCreateMeetingFromOrganization } from '../../../hooks/useMeeting';
import { TimelineToolbelt } from '../../ui-kit/molecules';

interface ToolbeltProps {
  organizationId: string;
}

export const OrginizationToolbelt: React.FC<ToolbeltProps> = ({
  organizationId,
}) => {
  const { onCreateOrganizationNote, saving } = useCreateOrganizationNote({
    organizationId,
  });
  const { onCreateMeeting } = useCreateMeetingFromOrganization({
    organizationId,
  });
  // const { onCreateMeeting } = useCreateMeetingFromOrganization({ organizationId });

  return (
    <TimelineToolbelt
      onCreateMeeting={onCreateMeeting}
      onCreateNote={onCreateOrganizationNote}
    />
  );
};
