import { NodeProps } from '@xyflow/react';
import { FlowActionType } from '@store/Flows/types';

import { Handle } from '../components';
import { EmailActionNode, SendConnectionRequestActionNode } from './actions';

interface BaseNodeData {
  isEditing?: boolean;
}

interface EmailData extends BaseNodeData {
  subject: string;
  bodyTemplate: string;
  action: FlowActionType.EMAIL_NEW | FlowActionType.EMAIL_REPLY;
}

interface LinkedInData extends BaseNodeData {
  messageTemplate: string;
  action:
    | FlowActionType.LINKEDIN_MESSAGE
    | FlowActionType.LINKEDIN_CONNECTION_REQUEST;
}

type ActionNodeData = EmailData | LinkedInData;

type ActionNodeProps = Omit<NodeProps, 'data'> & {
  data: ActionNodeData;
};

export const ActionNode = ({ data, ...flowProps }: ActionNodeProps) => {
  const renderNode = () => {
    if (isEmailData(data)) {
      return (
        <EmailActionNode
          {...flowProps}
          data={{
            isEditing: data.isEditing,
            subject: data.subject,
            bodyTemplate: data.bodyTemplate,
            action: data.action,
          }}
        />
      );
    }

    if (
      isLinkedInData(data) &&
      data.action === FlowActionType.LINKEDIN_CONNECTION_REQUEST
    ) {
      return (
        <SendConnectionRequestActionNode
          {...flowProps}
          data={{
            isEditing: data.isEditing,
            messageTemplate: data.messageTemplate,
            action: data.action,
          }}
        />
      );
    }

    return null;
  };

  return (
    <>
      <div className='h-[56px] max-w-[300px] w-[300px] bg-white border border-grayModern-300 p-4 rounded-lg group cursor-pointer flex items-center'>
        {renderNode()}
        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};

const isEmailData = (data: ActionNodeData): data is EmailData => {
  return (
    'subject' in data &&
    'bodyTemplate' in data &&
    (data.action === FlowActionType.EMAIL_NEW ||
      data.action === FlowActionType.EMAIL_REPLY)
  );
};

const isLinkedInData = (data: ActionNodeData): data is LinkedInData => {
  return (
    'messageTemplate' in data &&
    (data.action === FlowActionType.LINKEDIN_MESSAGE ||
      data.action === FlowActionType.LINKEDIN_CONNECTION_REQUEST)
  );
};
