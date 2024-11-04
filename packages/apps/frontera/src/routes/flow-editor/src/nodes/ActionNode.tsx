import { NodeProps } from '@xyflow/react';
import { FlowActionType } from '@store/Flows/types';

import { Handle } from '../components';
import { EmailActionNode, SendConnectionRequestActionNode } from './actions';

export const ActionNode = (
  props: NodeProps & {
    data: {
      subject: string;
      isEditing?: boolean;
      bodyTemplate: string;
      action: FlowActionType;
    };
  },
) => {
  const action = props.data.action;

  return (
    <>
      <div
        className={`h-[56px] max-w-[300px] w-[300px] bg-white border border-grayModern-300 p-4 rounded-lg group cursor-pointer flex items-center`}
      >
        {[FlowActionType.EMAIL_NEW, FlowActionType.EMAIL_REPLY].includes(
          action,
        ) && <EmailActionNode {...props} />}
        {FlowActionType.LINKEDIN_CONNECTION_REQUEST === action && (
          <SendConnectionRequestActionNode />
        )}

        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};
