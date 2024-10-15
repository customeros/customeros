import { Node } from '@xyflow/react';

import { EmailSettings } from './EmailSettings';
import { NoEmailNodesPanel } from './NoEmailNodesPanel';

export const FlowSettingsPanel = ({
  id,
  nodes,
}: {
  id: string;
  nodes: Node[];
}) => {
  const hasEmailNodes = nodes.some((node) => node.data?.action === 'EMAIL_NEW');

  return (
    <div className='absolute z-10 top-[41px] bottom-0 right-0 w-[400px] bg-white p-4 border-l flex flex-col gap-4 animate-slideLeft'>
      {hasEmailNodes && <EmailSettings id={id} />}
      {!hasEmailNodes && <NoEmailNodesPanel />}
    </div>
  );
};
