import { Action, InteractionEvent, Meeting } from '@graphql/types';

import { EmailPreviewModal } from '../events/email/EmailPreviewModal';
import { MeetingPreviewModal } from '../events/meeting/MeetingPreviewModal';
import { SlackThreadPreviewModal } from '../events/slack/SlackThreadPreviewModal';
import { ActionPreviewModal } from '../events/action/ActionPreviewModal';

import { useTimelineEventPreviewContext } from './TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { TimelinePreviewBackdrop } from '@organization/components/Timeline/preview/TimelinePreviewBackdrop';
import { IntercomThreadPreviewModal } from '@organization/components/Timeline/events/intercom/IntercomThreadPreviewModal';
import { LogEntryPreviewModal } from '@organization/components/Timeline/events/logEntry/LogEntryPreviewModal';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';

interface TimelineEventPreviewModalProps {
  invalidateQuery: () => void;
}

export const TimelineEventPreviewModal = ({
  invalidateQuery,
}: TimelineEventPreviewModalProps) => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();

  const event = modalContent as
    | InteractionEvent
    | Meeting
    | Action
    | LogEntryWithAliases;
  const isMeeting = event?.__typename === 'Meeting';
  const isAction = event?.__typename === 'Action';
  const isLogEntry = event?.__typename === 'LogEntry';
  const isInteraction = event?.__typename === 'InteractionEvent';
  const isSlack = isInteraction && event?.channel === 'SLACK';
  const isIntercom = isInteraction && event?.channel === 'CHAT';
  const isEmail = isInteraction && event?.channel === 'EMAIL';

  // Email handles close logic from within and use outside click cannot be used because preview should be closed only on backdrop click
  // user should be able to update panel details while having preview open
  if (isEmail) {
    return <EmailPreviewModal invalidateQuery={invalidateQuery} />;
  }

  return (
    <TimelinePreviewBackdrop onCloseModal={closeModal}>
      {isMeeting && <MeetingPreviewModal invalidateQuery={invalidateQuery} />}
      {isSlack && <SlackThreadPreviewModal />}
      {isIntercom && <IntercomThreadPreviewModal />}
      {isAction && <ActionPreviewModal type={event.actionType} />}
      {isLogEntry && <LogEntryPreviewModal />}
    </TimelinePreviewBackdrop>
  );
};
