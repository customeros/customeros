import { OrganizationRelationship } from '@graphql/types';

export type RelationshipType =
  | 'Customer'
  | 'Prospect'
  | 'Stranger'
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
    label: 'Stranger',
    value: OrganizationRelationship.Stranger,
  },
  {
    label: 'Former Customer',
    value: OrganizationRelationship.FormerCustomer,
  },
];
