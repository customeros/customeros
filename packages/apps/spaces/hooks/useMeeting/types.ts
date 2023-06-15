export {
  useCreateMeetingMutation,
  useUpdateMeetingMutation,
  GetOrganizationTimelineDocument,
  GetContactTimelineDocument,
  useMeetingLinkAttachmentMutation,
  useMeetingUnlinkAttachmentMutation,
  useMeetingLinkRecordingMutation,
  useMeetingUnlinkRecordingMutation,
  useLinkMeetingAttendeeMutation,
  useUnlinkMeetingAttendeeMutation,
} from '../../graphQL/__generated__/generated';
export type {
  MeetingInput,
  MeetingParticipant,
  CreateMeetingMutation,
  GetContactTimelineQuery,
  UpdateMeetingMutation,
  MeetingUnlinkAttachmentMutation,
  MeetingUnlinkRecordingMutation,
  Meeting,
  GetOrganizationTimelineQuery,
  LinkMeetingAttendeeMutation,
  UnlinkMeetingAttendeeMutation,
} from '../../graphQL/__generated__/generated';

export const NOW_DATE = new Date().toISOString();

export interface Result {
  onCreateMeeting: (userId: string) => void;
}
