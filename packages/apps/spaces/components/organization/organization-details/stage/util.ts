export type Stage =
  | 'Target'
  | 'Lead'
  | 'Prospect'
  | 'Trial'
  | 'Lost'
  | 'Live'
  | 'Former'
  | 'Unqualified';

export const stageOptions: { label: string; value: Stage }[] = [
  { label: 'Target', value: 'Target' },
  { label: 'Lead', value: 'Lead' },
  { label: 'Prospect', value: 'Prospect' },
  { label: 'Trial', value: 'Trial' },
  { label: 'Lost', value: 'Lost' },
  { label: 'Live', value: 'Live' },
  { label: 'Former', value: 'Former' },
  { label: 'Unqualified', value: 'Unqualified' },
];
