import { useRef, ReactElement } from 'react';
import { useParams } from 'react-router-dom';

import { Node } from '@xyflow/react';
import { observer } from 'mobx-react-lite';
import { format, toZonedTime } from 'date-fns-tz';
import { FlowActionType } from '@store/Flows/types';
import { FlowStore } from '@store/Flows/Flow.store';

import { Mail01 } from '@ui/media/icons/Mail01';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { MailReply } from '@ui/media/icons/MailReply';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

interface ExtendedNode extends Node {
  internalId: string;
}

const FLOW_ACTION_ICONS: Record<string, ReactElement> = {
  [FlowActionType.EMAIL_NEW]: <Mail01 className='size-3' />,
  [FlowActionType.EMAIL_REPLY]: <MailReply className='size-3' />,
  [FlowActionType.LINKEDIN_CONNECTION_REQUEST]: (
    <LinkedinOutline className='size-3' />
  ),
  [FlowActionType.LINKEDIN_MESSAGE]: <LinkedinOutline className='size-3' />,
};
const TOOLTIP_MESSAGE: Record<string, (date: string) => string> = {
  [FlowActionType.EMAIL_NEW]: (date) => `Email scheduled for ${date}`,
  [FlowActionType.EMAIL_REPLY]: (date) => `Email scheduled for ${date}`,
  [FlowActionType.LINKEDIN_CONNECTION_REQUEST]: (date) =>
    `Connection request scheduled for ${date}`,
  [FlowActionType.LINKEDIN_MESSAGE]: (date) =>
    `LinkedIn message scheduled for ${date}`,
};

export const NextFlowAction = observer(
  ({ contactID }: { contactID: string }) => {
    const { flows, contacts } = useStore();
    const { id } = useParams<{ id: string }>();
    const itemRef = useRef<HTMLDivElement>(null);

    // Early return if required data is missing
    if (!id || !flows.value.has(id)) {
      return <span className='text-grayModern-400'>None</span>;
    }

    const flowStore = flows.value.get(id) as FlowStore;
    const contact = flowStore?.value?.participants.find(
      (c) => c.entityId === contactID,
    );

    if (!contact?.executions?.length || !flowStore?.value?.nodes) {
      return <span className='text-grayModern-400'>None</span>;
    }

    const contactTimezone = contacts?.value?.get(contactID)?.value?.timezone;

    // Process flow data
    const nodes = flowStore.parsedNodes;

    const actionNodes = nodes.filter(
      (node: ExtendedNode) => node.type === 'action',
    );

    // console.log(actionNodes);

    const nextAction = contact.executions.find(
      (e) => e.scheduledAt && !e.executedAt,
    );

    const nextActionIndex = contact.executions.findIndex(
      (e) => e?.scheduledAt && !e?.executedAt,
    );

    if (!nextAction) {
      return <span className='text-grayModern-400'>None</span>;
    }

    const nextActionNode = nodes.find(
      (e: ExtendedNode) => e.internalId === nextAction?.action?.metadata?.id,
    );
    const nextActionDate = nextAction.scheduledAt;

    const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

    const userTimezoneScheduledAt = toZonedTime(nextActionDate, userTimezone);
    const contactTimezoneScheduledAt =
      contactTimezone && toZonedTime(nextActionDate, contactTimezone);

    const formattedDate = format(
      userTimezoneScheduledAt,
      'd MMM y, HH:mm zzz',
      {
        timeZone: userTimezone,
      },
    );

    const formattedDateContact =
      contactTimezone &&
      contactTimezoneScheduledAt &&
      format(contactTimezoneScheduledAt, 'HH:mm zzz', {
        timeZone: contactTimezone,
      });

    if (typeof nextActionNode?.data?.action !== 'string') {
      return null;
    }

    return (
      <Tooltip
        hasArrow
        align='start'
        side='bottom'
        label={
          <div className='space-y-1'>
            <div className='flex gap-1'>
              {TOOLTIP_MESSAGE[nextActionNode.data.action](formattedDate)}
              {formattedDateContact && <span>({formattedDateContact})</span>}
            </div>
          </div>
        }
      >
        <div ref={itemRef}>
          <div
            data-test='flow-name'
            className='flex items-center gap-2 px-1.5 bg-gray-100 rounded-md w-max overflow-hidden'
          >
            {FLOW_ACTION_ICONS[nextActionNode.data.action]}
            <div className='text-sm truncate'>
              Step {nextActionIndex + 1}/{actionNodes.length}
            </div>
          </div>
        </div>
      </Tooltip>
    );
  },
);
