import { FlowSequenceStatus } from '@graphql/types';

export const statusOptions: {
  label: string;
  value: FlowSequenceStatus;
}[] = [
  {
    label: 'Active',
    value: FlowSequenceStatus.Active,
  },
  {
    value: FlowSequenceStatus.Archived,
    label: 'Archived',
  },
  {
    label: 'Inactive',
    value: FlowSequenceStatus.Inactive,
  },
  {
    label: 'Paused',
    value: FlowSequenceStatus.Paused,
  },
];
