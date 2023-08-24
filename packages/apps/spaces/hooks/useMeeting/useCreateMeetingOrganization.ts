import {
  GetOrganizationTimelineQuery,
  GetOrganizationTimelineDocument,
  useCreateMeetingMutation,
  NOW_DATE,
  Result,
} from './types';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';
import { DataSource } from '@spaces/graphql';
export interface Props {
  organizationId?: string;
}
export const useCreateMeetingFromOrganization = ({
  organizationId,
}: Props): Result => {
  const [createMeetingMutation, { loading, error, data }] =
    useCreateMeetingMutation({
      onError: () => {
        toast.error(
          `Something went wrong while adding draft meeting to the timeline`,
        );
      },
    });

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
      createdAt: new Date(),
      meetingStartedAt: new Date(),
      meetingEndedAt: new Date(),
      agendaContentType: 'text/html',
      meetingCreatedBy: meeting_Create.createdBy,
      describedBy: [],
      includes: [],
      events: [],
      recording: null,
      source: DataSource.Openline,
      ...meeting_Create,
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
    return createMeetingMutation({
      variables: {
        meeting: {
          createdBy: [{ userId }],
          attendedBy: [],
          appSource: 'OPENLINE',
          name: '',
          startedAt: new Date().toISOString(),
          endedAt: new Date().toISOString(),
          agenda: `<p>INTRODUCTION</p>
                     <p>DISCUSSION</p>
                     <p>NEXT STEPS</p>
                     `,
          agendaContentType: 'text/html',
          note: { html: '<p>Notes:</p>', appSource: 'OPENLINE' },
        },
      },
      update: handleUpdateCacheAfterAddingMeeting,
    });
  };

  return {
    onCreateMeeting: handleCreateMeetingFromOrganization,
  };
};
