import { useMemo, ReactElement } from 'react';

import { htmlToText } from 'html-to-text';
import { FlowActionType } from '@store/Flows/types';
import { NodeProps, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { MailReply } from '@ui/media/icons/MailReply';

import { EmailEditorModal } from '../../components';

const iconMap: Record<string, ReactElement> = {
  [FlowActionType.EMAIL_NEW]: <Mail01 className='text-inherit' />,
  [FlowActionType.EMAIL_REPLY]: <MailReply className='text-inherit' />,
};

const colorMap: Record<string, string> = {
  [FlowActionType.EMAIL_NEW]: 'blue',
  [FlowActionType.EMAIL_REPLY]: 'blue',
};

export const EmailActionNode = ({
  id,
  data,
}: NodeProps & {
  data: {
    subject: string;
    isEditing?: boolean;
    bodyTemplate: string;
    action: FlowActionType;
  };
}) => {
  const { setNodes } = useReactFlow();

  const color = colorMap?.[data.action];

  const handleEmailDataChange = ({
    subject,
    bodyTemplate,
  }: {
    subject: string;
    bodyTemplate: string;
  }) => {
    setNodes((nds) =>
      nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              subject,
              bodyTemplate,
              isEditing: false,
            },
          };
        }

        if (node.data?.replyTo === id) {
          return {
            ...node,
            data: {
              ...node.data,
              subject: `RE: ${subject}`,
            },
          };
        }

        return node;
      }),
    );
  };

  const handleCancel = () => {
    setNodes((nds) =>
      nds.map((node) => {
        if (node.id === id) {
          return {
            ...node,
            data: {
              ...node.data,
              isEditing: false,
            },
          };
        }

        return node;
      }),
    );
  };

  const parsedTemplate = useMemo(
    () => htmlToText(data?.bodyTemplate).trim(),
    [data?.bodyTemplate],
  );

  const toggleEditing = () => {
    setNodes((nds) =>
      nds.map((node) =>
        node.id === id
          ? { ...node, data: { ...node.data, isEditing: true } }
          : node,
      ),
    );
  };

  return (
    <>
      <div className='text-sm flex items-center justify-between overflow-hidden w-full'>
        <div className='truncate text-sm flex items-center'>
          <div
            className={cn(
              `size-6 min-w-6 mr-2 bg-${color}-50 text-${color}-500 border border-gray-100 rounded flex items-center justify-center`,
              {
                'border-blue-100': color === 'blue',
              },
            )}
          >
            {iconMap?.[data.action]}
          </div>
          <span className='truncate font-medium'>
            {data.subject?.length > 0 ? (
              data.subject
            ) : parsedTemplate?.length > 0 ? (
              parsedTemplate
            ) : (
              <span className='text-gray-400 font-normal'>
                Write an email that wows them
              </span>
            )}
          </span>
        </div>

        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Edit'
          icon={<Edit03 />}
          onClick={toggleEditing}
          className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
        />
      </div>
      <EmailEditorModal
        data={data}
        handleCancel={handleCancel}
        isEditorOpen={data.isEditing || false}
        handleEmailDataChange={handleEmailDataChange}
      />
    </>
  );
};
