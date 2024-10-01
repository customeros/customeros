import { match } from 'ts-pattern';

import { FlowContactStatus } from '@graphql/types';

interface FlowStatusCellProps {
  value: string;
}

export const FlowStatusCell = ({ value }: FlowStatusCellProps) => {
  const flowStatus = match(value)
    .with(FlowContactStatus.InProgress, () => 'In Progress')
    .with(FlowContactStatus.Paused, () => 'Paused')
    .with(FlowContactStatus.Completed, () => 'Completed')
    .with(FlowContactStatus.Scheduled, () => 'Scheduled')
    // Temporary: this should be replced with correct enum values generated in graphql
    .with('PENDING', () => 'Pending')
    .with('GOAL_ACHIEVED', () => 'Goal achieved')
    .otherwise(() => value);

  return <div>{flowStatus}</div>;
};
