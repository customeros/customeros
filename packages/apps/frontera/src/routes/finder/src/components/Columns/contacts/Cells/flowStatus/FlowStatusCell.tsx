import { match } from 'ts-pattern';

import { FlowParticipantStatus } from '@graphql/types';

interface FlowStatusCellProps {
  value: string;
}

export const FlowStatusCell = ({ value }: FlowStatusCellProps) => {
  const flowStatus = match(value)
    .with(FlowParticipantStatus.Pending, () => 'Pending')
    .with(FlowParticipantStatus.InProgress, () => 'In Progress')
    .with(FlowParticipantStatus.Paused, () => 'Paused')
    .with(FlowParticipantStatus.Completed, () => 'Completed')
    .with(FlowParticipantStatus.Scheduled, () => 'Scheduled')
    .with(FlowParticipantStatus.GoalAchieved, () => 'Goal achieved')
    .otherwise(() => <span className='text-grayModern-400'>Not in flow</span>);

  return <div>{flowStatus}</div>;
};
