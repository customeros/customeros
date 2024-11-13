import { FlowStatus } from '@graphql/types';

export const flowOptions = [
  { label: 'Live', value: FlowStatus.On },
  { label: 'Stopped', value: FlowStatus.Off },
];
