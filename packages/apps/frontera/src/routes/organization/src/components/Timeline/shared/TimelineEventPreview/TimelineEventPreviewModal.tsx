import { VirtuosoHandle } from 'react-virtuoso';

import {
  InvoiceWithId,
  LogEntryWithAliases,
} from '@organization/components/Timeline/types';
import {
  Issue,
  Action,
  Meeting,
  InteractionEvent,
  ExternalSystemType,
} from '@graphql/types';
import { IssuePreviewModal } from '@organization/components/Timeline/PastZone/events/issue/IssuePreviewModal';
import { InvoicePreviewModal } from '@organization/components/Timeline/PastZone/events/invoice/InvoicePreviewModal';
import { LogEntryPreviewModal } from '@organization/components/Timeline/PastZone/events/logEntry/LogEntryPreviewModal';
import { TimelinePreviewBackdrop } from '@organization/components/Timeline/shared/TimelineEventPreview/TimelinePreviewBackdrop';
import { IntercomThreadPreviewModal } from '@organization/components/Timeline/PastZone/events/intercom/IntercomThreadPreviewModal';
import { LogEntryUpdateModalContextProvider } from '@organization/components/Timeline/PastZone/events/logEntry/context/LogEntryUpdateModalContext';
import { useTimelineEventPreviewStateContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { EmailPreviewModal } from '../../PastZone/events/email/EmailPreviewModal';
import { ActionPreviewModal } from '../../PastZone/events/action/ActionPreviewModal';
import { MeetingPreviewModal } from '../../PastZone/events/meeting/MeetingPreviewModal';
import { SlackThreadPreviewModal } from '../../PastZone/events/slack/SlackThreadPreviewModal';

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
    | InvoiceWithId
    | LogEntryWithAliases;
  const isMeeting = event?.__typename === 'Meeting';
  const isAction = event?.__typename === 'Action';
  const isLogEntry = event?.__typename === 'LogEntry';
  const isInteraction = event?.__typename === 'InteractionEvent';
  const isIssue = event?.__typename === 'Issue';
  const isInvoice = event?.__typename === 'Invoice';
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
        virtuosoRef={virtuosoRef}
        invalidateQuery={invalidateQuery}
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
        {isLogEntry && (
          <LogEntryPreviewModal invalidateQuery={invalidateQuery} />
        )}
        {isIssue && <IssuePreviewModal />}
        {isInvoice && <InvoicePreviewModal />}
      </TimelinePreviewBackdrop>
    </LogEntryUpdateModalContextProvider>
  );
};
