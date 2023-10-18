import {
  Action,
  InteractionEvent,
  Meeting,
  ExternalSystemType,
  Issue,
} from '@graphql/types';

import { EmailPreviewModal } from '../events/email/EmailPreviewModal';
import { MeetingPreviewModal } from '../events/meeting/MeetingPreviewModal';
import { SlackThreadPreviewModal } from '../events/slack/SlackThreadPreviewModal';
import { ActionPreviewModal } from '../events/action/ActionPreviewModal';

import { TimelinePreviewBackdrop } from '@organization/src/components/Timeline/preview/TimelinePreviewBackdrop';
import { IntercomThreadPreviewModal } from '@organization/src/components/Timeline/events/intercom/IntercomThreadPreviewModal';
import { LogEntryPreviewModal } from '@organization/src/components/Timeline/events/logEntry/LogEntryPreviewModal';
import { LogEntryWithAliases } from '@organization/src/components/Timeline/types';
import { LogEntryUpdateModalContextProvider } from '@organization/src/components/Timeline/events/logEntry/context/LogEntryUpdateModalContext';
import { VirtuosoHandle } from 'react-virtuoso';
import { IssuePreviewModal } from '@organization/src/components/Timeline/events/issue/IssuePreviewModal';
import { useTimelineEventPreviewStateContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface TimelineEventPreviewModalProps {
  invalidateQuery: () => void;
  virtuosoRef?: React.RefObject<VirtuosoHandle>;
}

export const TimelineEventPreviewModal = ({
  invalidateQuery,
  virtuosoRef,
}: TimelineEventPreviewModalProps) => {
  const { modalContent } = useTimelineEventPreviewStateContext();

  const event = modalContent as
    | InteractionEvent
    | Meeting
    | Action
    | Issue
    | LogEntryWithAliases;
  const isMeeting = event?.__typename === 'Meeting';
  const isAction = event?.__typename === 'Action';
  const isLogEntry = event?.__typename === 'LogEntry';
  const isInteraction = event?.__typename === 'InteractionEvent';
  const isIssue = event?.__typename === 'Issue';
  const isSlack =
    isInteraction &&
    event?.channel === 'CHAT' &&
    event?.externalLinks?.[0].type === ExternalSystemType.Slack;
  const isIntercom =
    isInteraction &&
    event?.channel === 'CHAT' &&
    event?.externalLinks?.[0].type === ExternalSystemType.Intercom;
  const isEmail = isInteraction && event?.channel === 'EMAIL';

  // Email handles close logic from within and use outside click cannot be used because preview should be closed only on backdrop click
  // user should be able to update panel details while having preview open
  if (isEmail) {
    return (
      <EmailPreviewModal
        invalidateQuery={invalidateQuery}
        virtuosoRef={virtuosoRef}
      />
    );
  }

  return (
    <LogEntryUpdateModalContextProvider>
      <TimelinePreviewBackdrop>
        {isMeeting && <MeetingPreviewModal invalidateQuery={invalidateQuery} />}
        {isSlack && <SlackThreadPreviewModal />}
        {isIntercom && <IntercomThreadPreviewModal />}
        {isAction && <ActionPreviewModal type={event.actionType} />}
        {isLogEntry && <LogEntryPreviewModal />}
        {isIssue && <IssuePreviewModal />}
      </TimelinePreviewBackdrop>
    </LogEntryUpdateModalContextProvider>
  );
};
