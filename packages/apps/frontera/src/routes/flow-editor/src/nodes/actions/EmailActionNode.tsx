import { useParams } from 'react-router-dom';
import { useMemo, ReactElement, MouseEventHandler } from 'react';

import { htmlToText } from 'html-to-text';
import { observer } from 'mobx-react-lite';
import { FlowActionType } from '@store/Flows/types';
import { NodeProps, useReactFlow } from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { FlowStatus } from '@graphql/types';
import { Mail01 } from '@ui/media/icons/Mail01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Trash01 } from '@ui/media/icons/Trash01';
import { MailReply } from '@ui/media/icons/MailReply';

const iconMap: Record<string, ReactElement> = {
  [FlowActionType.EMAIL_NEW]: <Mail01 className='text-inherit' />,
  [FlowActionType.EMAIL_REPLY]: <MailReply className='text-inherit' />,
};

const colorMap: Record<string, string> = {
  [FlowActionType.EMAIL_NEW]: 'blue',
  [FlowActionType.EMAIL_REPLY]: 'blue',
};

export const EmailActionNode = observer(
  ({
    data,
    id,
  }: NodeProps & {
    data: {
      subject: string;
      isEditing?: boolean;
      bodyTemplate: string;
      action: FlowActionType;
    };
  }) => {
    const color = colorMap?.[data.action];
    const { ui, flows } = useStore();
    const flowId = useParams()?.id as string;
    const flow = flows.value.get(flowId)?.value;
    const flowWasStarted =
      flow?.status === FlowStatus.On || flow?.firstStartedAt;

    const { deleteElements } = useReactFlow();

    const parsedTemplate = useMemo(
      () => htmlToText(data?.bodyTemplate ?? '').trim(),
      [data?.bodyTemplate],
    );

    const handleDelete: MouseEventHandler = (e) => {
      e.stopPropagation();

      return deleteElements({ nodes: [{ id }] });
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

          {!ui.flowActionSidePanel.isOpen && !flowWasStarted && (
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='Edit'
              icon={<Trash01 />}
              onClick={handleDelete}
              className={`ml-2 opacity-0 group-hover:opacity-100 pointer-events-all`}
            />
          )}
        </div>
      </>
    );
  },
);
