export type Stage = 'ACTIVE' | 'INACTIVE' | 'PENDING' | 'REJECTED';

export const stageOptions: { label: string; value: Stage }[] = [
  { label: 'Active', value: 'ACTIVE' },
  { label: 'Inactive', value: 'INACTIVE' },
  { label: 'Pending', value: 'PENDING' },
  { label: 'Rejected', value: 'REJECTED' },
];
