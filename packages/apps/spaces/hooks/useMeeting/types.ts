export {
  useCreateMeetingMutation,
  useUpdateMeetingMutation,
  GetOrganizationTimelineDocument,
  GetContactTimelineDocument,
} from '../../graphQL/__generated__/generated';
export type {
  MeetingInput,
  MeetingParticipant,
  CreateMeetingMutation,
  GetContactTimelineQuery,
  UpdateMeetingMutation,
} from '../../graphQL/__generated__/generated';
import type {
  MeetingInput,
  MeetingParticipant,
  CreateMeetingMutation,
} from '../../graphQL/__generated__/generated';

export const NOW_DATE = new Date().toISOString();

export interface Result {
  onCreateMeeting: () => Promise<
    CreateMeetingMutation['meeting_Create'] | null
  >;
}
