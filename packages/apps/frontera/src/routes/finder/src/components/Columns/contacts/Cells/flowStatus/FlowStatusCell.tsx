import { useParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Rocket02 } from '@ui/media/icons/Rocket02';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { FlowParticipantStatus } from '@graphql/types';
import { Hourglass02 } from '@ui/media/icons/Hourglass02';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { CalendarCheck01 } from '@ui/media/icons/CalendarCheck01';

interface FlowStatusCellProps {
  contactID: string;
}

export const FlowStatusCell = observer(({ contactID }: FlowStatusCellProps) => {
  const { flows } = useStore();
  const id = useParams()?.id as string;

  const flowStore = flows.value.get(id)?.value;
  const contact = flowStore?.participants.find((c) => c.entityId === contactID);

  const flowStatus = match(contact?.status)
    .with(FlowParticipantStatus.OnHold, () => [
      'Blocked',
      <SlashCircle01 className='size-3' />,
    ])
    .with(FlowParticipantStatus.Ready, () => [
      'Ready',
      <Rocket02 className='size-3' />,
    ])
    .with(FlowParticipantStatus.InProgress, () => [
      'In Progress',
      <Hourglass02 className='size-3' />,
    ])
    .with(FlowParticipantStatus.Completed, () => [
      'Completed',
      <CheckCircle className='size-3' />,
    ])
    .with(FlowParticipantStatus.Scheduled, () => [
      'Scheduled',
      <CalendarCheck01 className='size-3' />,
    ])
    .with(FlowParticipantStatus.GoalAchieved, () => [
      'Goal achieved',
      <Trophy01 className='size-3' />,
    ])
    .otherwise(() => [
      <span className='text-grayModern-400'>Not in flow</span>,
      null,
    ]);

  return (
    <div className='flex items-center overflow-x-hidden gap-2 overflow-ellipsis bg-gray-100 rounded-md px-1.5 truncate w-[fit-content]'>
      {flowStatus?.[1] && <span className='flex'>{flowStatus[1]}</span>}
      <div className='truncate'>{flowStatus[0]}</div>
    </div>
  );
});
