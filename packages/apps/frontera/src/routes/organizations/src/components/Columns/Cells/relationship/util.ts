export type RelationshipType = 'Customer' | 'Prospect';

export const relationshipOptions: {
  value: boolean;
  label: RelationshipType;
}[] = [
  {
    value: true,
    label: 'Customer',
  },
  {
    value: false,
    label: 'Prospect',
  },
];
