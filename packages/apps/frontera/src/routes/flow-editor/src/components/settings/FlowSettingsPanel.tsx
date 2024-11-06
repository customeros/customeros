import { Node } from '@xyflow/react';
import { FlowActionType } from '@store/Flows/types.ts';

import { SenderSettings } from './SenderSettings.tsx';
import { NoEmailNodesPanel } from './NoEmailNodesPanel';

export const FlowSettingsPanel = ({
  id,
  nodes,
}: {
  id: string;
  nodes: Node[];
}) => {
  const hasEmailNodes = nodes.some(
    (node) =>
      node.data?.action &&
      [FlowActionType.EMAIL_NEW, FlowActionType.EMAIL_REPLY].includes(
        node.data.action as FlowActionType,
      ),
  );
  const hasLinkedInNodes = nodes.some(
    (node) =>
      node.data?.action &&
      [
        FlowActionType.LINKEDIN_CONNECTION_REQUEST,
        FlowActionType.LINKEDIN_MESSAGE,
      ].includes(node.data.action as FlowActionType),
  );

  const showSenderSettings = hasEmailNodes || hasLinkedInNodes;

  return (
    <div className='absolute z-10 top-[41px] bottom-0 right-0 w-[400px] bg-white p-4 border-l flex flex-col gap-4 animate-slideLeft'>
      {showSenderSettings && (
        <SenderSettings
          id={id}
          hasEmailNodes={hasEmailNodes}
          hasLinkedInNodes={hasLinkedInNodes}
        />
      )}
      {!showSenderSettings && <NoEmailNodesPanel />}
    </div>
  );
};
