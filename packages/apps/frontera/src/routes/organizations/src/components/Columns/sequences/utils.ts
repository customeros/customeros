import { FlowSequenceStatus } from '@graphql/types';

export const flowSequencesOptions = [
  { label: 'Live', value: FlowSequenceStatus.Active },
  { label: 'Not Started', value: FlowSequenceStatus.Inactive },
  { label: 'Paused', value: FlowSequenceStatus.Paused },
];
