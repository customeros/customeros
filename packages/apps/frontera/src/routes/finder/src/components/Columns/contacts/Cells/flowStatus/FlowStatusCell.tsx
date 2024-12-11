import { ReactElement } from 'react';
import { useParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Rocket02 } from '@ui/media/icons/Rocket02';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { RefreshCw02 } from '@ui/media/icons/RefreshCw02';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { CalendarCheck01 } from '@ui/media/icons/CalendarCheck01';
import {
  FlowParticipantStatus,
  FlowParticipantRequirementsUnmeet,
} from '@graphql/types';

interface FlowStatusCellProps {
  contactID: string;
}

const requirementsCopy = {
  [FlowParticipantRequirementsUnmeet.MissingLinkedinUrl]: 'Missing LinkedIn',
  [FlowParticipantRequirementsUnmeet.MissingPrimaryEmail]: 'Missing email',
};

export const FlowStatusCell = observer(({ contactID }: FlowStatusCellProps) => {
  const { flows } = useStore();
  const id = useParams()?.id as string;

  const flowStore = flows.value.get(id)?.value;
  const contact = flowStore?.participants.find((c) => c.entityId === contactID);

  return match(contact?.status)
    .with(FlowParticipantStatus.OnHold, () => {
      if (!contact?.requirementsUnmeet?.length) {
        return (
          <StatusTag
            status='Blocked'
            icon={<SlashCircle01 className='size-3' />}
          />
        );
      }
      const copy =
        contact?.requirementsUnmeet?.length &&
        contact?.requirementsUnmeet?.length > 1
          ? 'Missing Details'
          : requirementsCopy[contact.requirementsUnmeet[0]];

      return (
        <Tooltip
          align='center'
          label={
            contact?.requirementsUnmeet?.length > 1 &&
            contact?.requirementsUnmeet.map((e) => <p>{requirementsCopy[e]}</p>)
          }
        >
          <div>
            <StatusTag
              status={copy}
              icon={<SlashCircle01 className='size-3' />}
            />
          </div>
        </Tooltip>
      );
    })
    .with(FlowParticipantStatus.Ready, () => (
      <StatusTag status='Ready' icon={<Rocket02 className='size-3' />} />
    ))
    .with(FlowParticipantStatus.InProgress, () => (
      <StatusTag
        status='In Progress'
        icon={<Hourglass02 className='size-3' />}
      />
    ))
    .with(FlowParticipantStatus.Completed, () => (
      <StatusTag status='Completed' icon={<CheckCircle className='size-3' />} />
    ))
    .with(FlowParticipantStatus.Scheduled, () => (
      <StatusTag
        status='Scheduled'
        icon={<CalendarCheck01 className='size-3' />}
      />
    ))
    .with(FlowParticipantStatus.Scheduling, () => (
      <StatusTag
        status='Scheduling'
        icon={<RefreshCw02 className='size-3' />}
      />
    ))
    .with(FlowParticipantStatus.GoalAchieved, () => (
      <StatusTag
        status='Goal achieved'
        icon={<Trophy01 className='size-3' />}
      />
    ))
    .otherwise(() => [
      <span className='text-grayModern-400'>No status yet</span>,
      null,
    ]);
});

const StatusTag = ({
  status,
  icon,
}: {
  status: string;
  icon?: ReactElement;
}) => {
  return (
    <div className='flex items-center overflow-x-hidden gap-2 overflow-ellipsis bg-gray-100 rounded-md px-1.5 truncate w-[fit-content]'>
      {icon && <span className='flex'>{icon}</span>}
      <div className='truncate'>{status}</div>
    </div>
  );
};
