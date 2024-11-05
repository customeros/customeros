import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Rocket02 } from '@ui/media/icons/Rocket02';
import { FlowParticipantStatus } from '@graphql/types';
import { TableCellTooltip } from '@ui/presentation/Table';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { CalendarCheck01 } from '@ui/media/icons/CalendarCheck01';

interface ContactNameCellProps {
  contactId: string;
}

const icons = {
  [FlowParticipantStatus.OnHold]: <SlashCircle01 className='size-3' />,
  [FlowParticipantStatus.Ready]: <Rocket02 className='size-3' />,
  [FlowParticipantStatus.Scheduled]: <CalendarCheck01 className='size-3' />,
  [FlowParticipantStatus.InProgress]: <Hourglass02 className='size-3' />,
  [FlowParticipantStatus.Completed]: <CheckCircle className='size-3' />,
  [FlowParticipantStatus.GoalAchieved]: <Trophy01 className='size-3' />,
};

export const ContactFlowCell = observer(
  ({ contactId }: ContactNameCellProps) => {
    const store = useStore();
    const [isEditing, setIsEditing] = useState(false);

    const contactStore = store.contacts.value.get(contactId);
    const itemRef = useRef<HTMLDivElement>(null);

    const contactFlows = contactStore?.flows;

    const open = () => {
      setIsEditing(true);
      store.ui.commandMenu.setType('EditContactFlow');
      store.ui.commandMenu.setOpen(true);
    };

    if (!contactFlows?.length && !isEditing) {
      return (
        <div
          onDoubleClick={open}
          className={cn(
            'flex w-full gap-1 items-center [&_.edit-button]:hover:opacity-100',
          )}
        >
          <div className='text-gray-400'>None</div>
          <IconButton
            size='xxs'
            onClick={open}
            variant='ghost'
            id='edit-button'
            aria-label='edit owner'
            className='edit-button opacity-0'
            dataTest={`contact-flow-edit-${contactId}`}
            icon={<Edit03 className='text-gray-500 size-3' />}
          />
        </div>
      );
    }

    const status = contactFlows?.[0]?.value.contacts.find(
      (e) => e.contact.metadata.id === contactId,
    )?.status;

    return (
      <TableCellTooltip
        hasArrow
        align='start'
        side='bottom'
        targetRef={itemRef}
        label={
          <div>
            {contactFlows?.map((flow) => (
              <div className='flex gap-1' key={flow.value.metadata.id}>
                <div>
                  {flow.value.name} •{' '}
                  <span className='capitalize'>
                    {flow.value.contacts
                      .find((e) => e.contact.metadata.id === contactId)
                      ?.status?.toLowerCase()
                      ?.split('_')
                      ?.join(' ')}{' '}
                  </span>
                </div>
              </div>
            ))}
          </div>
        }
      >
        <div
          onDoubleClick={open}
          className={cn(
            'cursor-default overflow-hidden overflow-ellipsis flex gap-1  [&_.edit-button]:hover:opacity-100',
          )}
        >
          <div ref={itemRef} className='flex overflow-hidden'>
            <div
              data-test='flow-name'
              className='flex items-center overflow-x-hidden gap-2 overflow-ellipsis bg-gray-100 rounded-md px-1.5 truncate'
            >
              {status && icons?.[status]}
              <div className='truncate'>{contactFlows?.[0]?.value.name}</div>
            </div>{' '}
            {!!contactFlows?.length && contactFlows.length > 1 && (
              <div className='rounded-md w-fit px-1.5 ml-1 text-gray-500 '>
                +{contactFlows?.length - 1}
              </div>
            )}
          </div>
          <IconButton
            size='xxs'
            onClick={open}
            variant='ghost'
            id='edit-button'
            aria-label='edit owner'
            className='edit-button opacity-0'
            icon={<Edit03 className='text-gray-500 size-3' />}
          />
        </div>
      </TableCellTooltip>
    );
  },
);
