import { useRef, ReactElement } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
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

const icons: Record<string, ReactElement> = {
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

    const contactStore = store.contacts.value.get(contactId);
    const itemRef = useRef<HTMLDivElement>(null);

    const contactFlows = contactStore?.flows;

    const open = () => {
      store.ui.commandMenu.setType('EditContactFlow');
      store.ui.commandMenu.setOpen(true);
    };

    if (!contactFlows?.length) {
      return (
        <div
          onClick={open}
          className={cn(
            'flex w-full gap-1 items-center [&_.edit-button]:hover:opacity-100 cursor-pointer',
          )}
        >
          <div
            className='text-gray-400'
            data-test={`contact-current-flows-in-contacts-table`}
          >
            None
          </div>
        </div>
      );
    }

    const status = contactFlows?.[0]?.value?.participants.find(
      (e) => e.entityId === contactId,
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
              <div className='flex gap-1' key={flow?.value?.metadata?.id}>
                <div>
                  {flow?.value?.name} •{' '}
                  <span className='capitalize'>
                    {flow?.value?.participants
                      .find((e) => e.entityId === contactId)
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
          onClick={open}
          className={cn(
            'overflow-hidden overflow-ellipsis flex gap-1 [&_.edit-button]:hover:opacity-100 cursor-pointer',
          )}
        >
          <div ref={itemRef} className='flex overflow-hidden'>
            <div
              data-test='flow-name'
              className='flex items-center overflow-x-hidden gap-2 overflow-ellipsis bg-gray-100 rounded-md px-1.5 truncate'
            >
              <span className='flex'>{status && icons?.[status]}</span>
              <div className='truncate'>{contactFlows?.[0]?.value.name}</div>
            </div>{' '}
            {!!contactFlows?.length && contactFlows.length > 1 && (
              <div className='rounded-md w-fit px-1.5 ml-1 text-gray-500 '>
                +{contactFlows?.length - 1}
              </div>
            )}
          </div>
        </div>
      </TableCellTooltip>
    );
  },
);
