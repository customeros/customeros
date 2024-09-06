import { FlowSequenceStatus } from '@graphql/types';

export const statusOptions: {
  label: string;
  value: FlowSequenceStatus;
}[] = [
  {
    label: 'Live',
    value: FlowSequenceStatus.Active,
  },

  {
    label: 'Not Started',
    value: FlowSequenceStatus.Inactive,
  },
  {
    label: 'Paused',
    value: FlowSequenceStatus.Paused,
  },
  {
    value: FlowSequenceStatus.Archived,
    label: 'Archived',
  },
];
