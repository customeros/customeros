import {
  GetOrganizationTimelineQuery,
  GetOrganizationTimelineDocument,
  useCreateMeetingMutation,
  NOW_DATE,
  Result,
} from './types';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';
export interface Props {
  organizationId?: string;
}
export const useCreateMeetingFromOrganization = ({
  organizationId,
}: Props): Result => {
  const [createMeetingMutation, { loading, error, data }] =
    useCreateMeetingMutation();

  const handleUpdateCacheAfterAddingMeeting = (
    cache: ApolloCache<any>,
    { data: { meeting_Create } }: any,
  ) => {
    const data: GetOrganizationTimelineQuery | null = client.readQuery({
      query: GetOrganizationTimelineDocument,
      variables: {
        organizationId,
        from: NOW_DATE,
        size: 10,
      },
    });

    const newMeeting = {
      ...meeting_Create,
      createdAt: new Date(),
      agenda: '',
      agendaContentType: 'text/html',
      meetingCreatedBy: meeting_Create.createdBy,
    };

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationTimelineDocument,
        data: {
          organization: {
            organizationId,
            timelineEvents: [newMeeting],
          },
          variables: { organizationId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        timelineEvents: [newMeeting],
      },
    };

    client.writeQuery({
      query: GetOrganizationTimelineDocument,
      data: newData,
      variables: {
        organizationId,
        from: NOW_DATE,
        size: 10,
      },
    });
  };

  const handleCreateMeetingFromOrganization: Result['onCreateMeeting'] = async (
    userId,
  ) => {
    try {
      const response = await createMeetingMutation({
        variables: {
          meeting: {
            createdBy: [{ userID: userId, type: 'user' }],
            attendedBy: [],
            appSource: 'OPENLINE',
            name: '',
            start: new Date().toISOString(),
            end: new Date().toISOString(),
          },
        },
        //@ts-expect-error fixme
        update: handleUpdateCacheAfterAddingMeeting,
      });

      toast.success(`Added draft meeting to the timeline`);
      return response.data?.meeting_Create ?? null;
    } catch (err) {
      console.error(err);
      toast.error(
        `Something went wrong while adding draft meeting to the timeline`,
      );
      return null;
    }
  };

  return {
    onCreateMeeting: handleCreateMeetingFromOrganization,
  };
};
