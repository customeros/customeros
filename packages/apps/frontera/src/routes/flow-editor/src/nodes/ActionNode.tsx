import { useMemo, useState, ReactElement } from 'react';

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
    bodyTemplate: string;
    action: FlowActionType;
  };
}) => {
  const [isEditorOpen, setIsEditorOpen] = useState(false);
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

  return (
    <>
      <div
        className={`aspect-[9/1] max-w-[300px] bg-white border border-grayModern-300 p-3 rounded-lg group cursor-pointer`}
      >
        <div className='truncate  text-sm flex items-center '>
          <div className='truncate text-sm flex items-center'>
            <div
              className={`size-6 mr-2 bg-${color}-50 text-${color}-500 border border-${color}-100  rounded flex items-center justify-center`}
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
            size='xs'
            variant='ghost'
            aria-label='Edit'
            icon={<Edit03 />}
            onClick={() => setIsEditorOpen(true)}
            className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
          />
        </div>
        <EmailEditorModal
          data={data}
          isEditorOpen={isEditorOpen}
          setIsEditorOpen={setIsEditorOpen}
          handleEmailDataChange={handleEmailDataChange}
        />
        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};
