import { NodeProps } from '@xyflow/react';
import { FlowActionType } from '@store/Flows/types';

import { Handle } from '../components';
import { EmailActionNode } from './actions';

interface BaseNodeData {
  isEditing?: boolean;
}

interface EmailData extends BaseNodeData {
  subject: string;
  bodyTemplate: string;
  action: FlowActionType.EMAIL_NEW | FlowActionType.EMAIL_REPLY;
}

type ActionNodeData = EmailData;

type ActionNodeProps = Omit<NodeProps, 'data'> & {
  data: ActionNodeData;
};

export const ActionNode = ({ data, ...flowProps }: ActionNodeProps) => {
  return (
    <>
      <div className='h-[56px] max-w-[300px] w-[300px] bg-white border border-grayModern-300 p-4 rounded-lg group cursor-pointer flex items-center'>
        <EmailActionNode
          {...flowProps}
          data={{
            isEditing: data.isEditing,
            subject: data.subject,
            bodyTemplate: data.bodyTemplate,
            action: data.action,
          }}
        />
        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};
