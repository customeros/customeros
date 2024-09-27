import { OrganizationRelationship } from '@graphql/types';

export type RelationshipType =
  | 'Customer'
  | 'Prospect'
  | 'Not a Fit'
  | 'Former Customer';

export const relationshipOptions: {
  label: RelationshipType;
  value: OrganizationRelationship;
}[] = [
  {
    label: 'Customer',
    value: OrganizationRelationship.Customer,
  },
  {
    label: 'Prospect',
    value: OrganizationRelationship.Prospect,
  },
  {
    label: 'Not a Fit',
    value: OrganizationRelationship.NotAFit,
  },
  {
    label: 'Former Customer',
    value: OrganizationRelationship.FormerCustomer,
  },
];
