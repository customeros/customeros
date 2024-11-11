import { FlowStatus } from '@graphql/types';

export const flowOptions = [
  { label: 'Live', value: FlowStatus.Active },
  { label: 'Not Started', value: FlowStatus.Inactive },
  { label: 'Paused', value: FlowStatus.Paused },
  { label: 'Scheduling', value: FlowStatus.Scheduling },
];
