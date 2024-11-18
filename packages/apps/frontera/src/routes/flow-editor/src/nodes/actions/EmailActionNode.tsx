import { useMemo, ReactElement } from 'react';

import { NodeProps } from '@xyflow/react';
import { htmlToText } from 'html-to-text';
import { FlowActionType } from '@store/Flows/types';

import { cn } from '@ui/utils/cn';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { MailReply } from '@ui/media/icons/MailReply';

const iconMap: Record<string, ReactElement> = {
  [FlowActionType.EMAIL_NEW]: <Mail01 className='text-inherit' />,
  [FlowActionType.EMAIL_REPLY]: <MailReply className='text-inherit' />,
};

const colorMap: Record<string, string> = {
  [FlowActionType.EMAIL_NEW]: 'blue',
  [FlowActionType.EMAIL_REPLY]: 'blue',
};

export const EmailActionNode = ({
  data,
}: NodeProps & {
  data: {
    subject: string;
    isEditing?: boolean;
    bodyTemplate: string;
    action: FlowActionType;
  };
}) => {
  const color = colorMap?.[data.action];

  const parsedTemplate = useMemo(
    () => htmlToText(data?.bodyTemplate).trim(),
    [data?.bodyTemplate],
  );

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

        {/* this is just visual guidance, clicking on the whole node performs actual action*/}
        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='Edit'
          icon={<Edit03 />}
          className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
        />
      </div>
    </>
  );
};
