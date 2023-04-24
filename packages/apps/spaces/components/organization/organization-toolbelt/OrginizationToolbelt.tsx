import React from 'react';

import { useCreateOrganizationNote } from '../../../hooks/useNote';
import { useCreateMeetingFromOrganization } from '../../../hooks/useMeeting';
import { TimelineToolbelt } from '../../ui-kit/molecules';
import { useRecoilState, useRecoilValue } from 'recoil';
import { contactNewItemsToEdit, userData } from '../../../state';
import { useUser } from '../../../hooks/useUser';

interface ToolbeltProps {
  organizationId: string;
}

export const OrginizationToolbelt: React.FC<ToolbeltProps> = ({
  organizationId,
}) => {
  const { identity: userEmail } = useRecoilValue(userData);
  const { data, loading, error } = useUser({ email: userEmail });
  const [itemsInEditMode, setItemToEditMode] = useRecoilState(
    contactNewItemsToEdit,
  );
  const { onCreateOrganizationNote, saving } = useCreateOrganizationNote({
    organizationId,
  });
  const { onCreateMeeting } = useCreateMeetingFromOrganization({
    organizationId,
  });

  const handleCreateNote = (data: any) =>
    onCreateOrganizationNote(data).then((response) => {
      if (response?.id) {
        setItemToEditMode({
          timelineEvents: [
            ...itemsInEditMode.timelineEvents,
            { id: response.id },
          ],
        });
      }
    });

  const handleCreateMeeting = () =>
    onCreateMeeting(data?.id).then((response) => {
      if (response?.id) {
        setItemToEditMode({
          timelineEvents: [
            ...itemsInEditMode.timelineEvents,
            { id: response.id },
          ],
        });
      }
    });
  return (
    <TimelineToolbelt
      onCreateMeeting={handleCreateMeeting}
      onCreateNote={handleCreateNote}
      isSkewed
    />
  );
};
