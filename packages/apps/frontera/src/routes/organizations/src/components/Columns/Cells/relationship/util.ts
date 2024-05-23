import { OrganizationRelationship } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';

export const relationshipOptions: SelectOption<OrganizationRelationship>[] = [
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
