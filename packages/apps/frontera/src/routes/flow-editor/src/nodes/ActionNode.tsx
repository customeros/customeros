import { useMemo, ReactElement } from 'react';

import { htmlToText } from 'html-to-text';
import { FlowActionType } from '@store/Flows/types.ts';
import { NodeProps, useReactFlow } from '@xyflow/react';

import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { MailReply } from '@ui/media/icons/MailReply';

import { Handle, EmailEditorModal } from '../components';

const iconMap: Record<string, ReactElement> = {
  [FlowActionType.EMAIL_NEW]: <Mail01 className='text-inherit' />,
  [FlowActionType.EMAIL_REPLY]: <MailReply className='text-inherit' />,
};

const colorMap: Record<string, string> = {
  [FlowActionType.EMAIL_NEW]: 'blue',
  [FlowActionType.EMAIL_REPLY]: 'blue',
};

export const ActionNode = ({
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
      <div
        className={`aspect-[9/1] max-w-[300px] w-[300px] bg-white border border-grayModern-300 p-3 rounded-lg group cursor-pointer`}
      >
        <div className='text-sm flex items-center justify-between '>
          <div className='truncate text-sm flex items-center'>
            <div
              className={`size-6 min-w-6 mr-2 bg-${color}-50 text-${color}-500 border border-${color}-100  rounded flex items-center justify-center`}
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
          isEditorOpen={data.isEditing || false}
          handleEmailDataChange={handleEmailDataChange}
          setIsEditorOpen={(isOpen: boolean) => {
            setNodes((nds) =>
              nds.map((node) =>
                node.id === id
                  ? { ...node, data: { ...node.data, isEditing: isOpen } }
                  : node,
              ),
            );
          }}
        />
        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};
