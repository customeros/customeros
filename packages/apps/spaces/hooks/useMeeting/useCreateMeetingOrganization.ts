import {
  GetContactTimelineQuery,
  GetOrganizationTimelineDocument,
  useCreateMeetingMutation,
  NOW_DATE,
  Result,
} from './types';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';
export interface Props {
  organizationId?: string;
}
export const useCreateMeetingFromOrganization = ({
  organizationId,
}: Props): Result => {
  const [createMeetingMutation, { loading, error, data }] =
    useCreateMeetingMutation();
  const loggedInUserData = useRecoilValue(userData);

  const handleUpdateCacheAfterAddingMeeting = (
    cache: ApolloCache<any>,
    { data: { meeting_Create } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetOrganizationTimelineDocument,
      variables: {
        organizationId,
        from: NOW_DATE,
        size: 10,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationTimelineDocument,
        data: {
          contact: {
            organizationId,
            timelineEvents: [meeting_Create],
          },
          variables: { organizationId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [
          ...(data.contact?.timelineEvents || []),
          meeting_Create,
        ],
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

  const handleCreateMeetingFromContact: Result['onCreateMeeting'] =
    async () => {
      try {
        const response = await createMeetingMutation({
          variables: {
            meeting: {
              createdBy: [{ userID: loggedInUserData.id, type: 'user' }],
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
    onCreateMeeting: handleCreateMeetingFromContact,
  };
};
