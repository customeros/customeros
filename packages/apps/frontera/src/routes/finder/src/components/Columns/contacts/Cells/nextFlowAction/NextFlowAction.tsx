import { useRef, ReactElement } from 'react';
import { useParams } from 'react-router-dom';

import { Node } from '@xyflow/react';
import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types';

import { Mail01 } from '@ui/media/icons/Mail01';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

interface ExtendedNode extends Node {
  internalId: string;
}

const FLOW_ACTION_ICONS: Record<string, ReactElement> = {
  [FlowActionType.EMAIL_NEW]: <Mail01 className='size-3' />,
  [FlowActionType.EMAIL_REPLY]: <Mail01 className='size-3' />,
  [FlowActionType.LINKEDIN_CONNECTION_REQUEST]: (
    <LinkedinOutline className='size-3' />
  ),
  [FlowActionType.LINKEDIN_MESSAGE]: <LinkedinOutline className='size-3' />,
};

const parseNodes = (nodesString: string): ExtendedNode[] => {
  try {
    return JSON.parse(nodesString);
  } catch (error) {
    console.error('Failed to parse flow nodes:', error);

    return [];
  }
};

const getActionNodes = (nodes: ExtendedNode[]): ExtendedNode[] =>
  nodes.filter((node) => node.type === 'action');

const formatActionName = (action: string): string =>
  action.split('_').join(' ').toLowerCase();

export const NextFlowAction = observer(
  ({ contactID }: { contactID: string }) => {
    const { flows } = useStore();
    const { id } = useParams<{ id: string }>();
    const itemRef = useRef<HTMLDivElement>(null);

    // Early return if required data is missing
    if (!id || !flows.value.has(id)) {
      return <span className='text-grayModern-400'>None</span>;
    }

    const flowStore = flows.value.get(id)?.value;
    const contact = flowStore?.participants.find(
      (c) => c.entityId === contactID,
    );

    if (!contact?.executions?.length || !flowStore?.nodes) {
      return <span className='text-grayModern-400'>None</span>;
    }

    // Process flow data
    const nodes = parseNodes(flowStore.nodes);
    const actionNodes = getActionNodes(nodes);

    const nextActionId = contact.executions.find(
      (e) => e.scheduledAt && !e.executedAt,
    )?.action.metadata.id;

    const nextActionNode = nodes.find((e) => e.internalId === nextActionId);
    const nextActionIndex = actionNodes.findIndex(
      (e) => e.internalId === nextActionId,
    );

    return (
      <Tooltip
        hasArrow
        align='start'
        side='bottom'
        label={
          <div className='space-y-1'>
            {actionNodes.map((action, index) => (
              <div className='flex gap-1' key={action.internalId}>
                <span>
                  Step {index + 1} • (
                  <span className='capitalize'>
                    {formatActionName(action.data.action as string)}
                  </span>
                  )
                </span>
              </div>
            ))}
          </div>
        }
      >
        <div ref={itemRef}>
          <div
            data-test='flow-name'
            className='flex items-center gap-2 px-1.5 bg-gray-100 rounded-md w-max overflow-hidden'
          >
            {typeof nextActionNode?.data?.action === 'string' &&
              FLOW_ACTION_ICONS[nextActionNode.data.action]}
            <div className='text-sm truncate'>
              Step {nextActionIndex + 1}/{actionNodes.length}
            </div>
          </div>
        </div>
      </Tooltip>
    );
  },
);
