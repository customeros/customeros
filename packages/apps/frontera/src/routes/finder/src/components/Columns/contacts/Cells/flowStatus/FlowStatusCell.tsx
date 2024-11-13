import { useParams } from 'react-router-dom';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { FlowParticipantStatus } from '@graphql/types';

interface FlowStatusCellProps {
  contactID: string;
}

export const FlowStatusCell = observer(({ contactID }: FlowStatusCellProps) => {
  const { flows } = useStore();
  const id = useParams()?.id as string;

  const flowStore = flows.value.get(id)?.value;
  const contact = flowStore?.contacts.find(
    (c) => c.contact.metadata.id === contactID,
  );

  const flowStatus = match(contact?.status)
    .with(FlowParticipantStatus.OnHold, () => 'Blocked')
    .with(FlowParticipantStatus.Ready, () => 'Ready')
    .with(FlowParticipantStatus.InProgress, () => 'In Progress')
    .with(FlowParticipantStatus.Completed, () => 'Completed')
    .with(FlowParticipantStatus.Scheduled, () => 'Scheduled')
    .with(FlowParticipantStatus.GoalAchieved, () => 'Goal achieved')
    .otherwise(() => <span className='text-grayModern-400'>Not in flow</span>);

  return <div>{flowStatus}</div>;
});
