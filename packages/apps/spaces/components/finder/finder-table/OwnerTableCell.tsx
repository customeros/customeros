import React from 'react';
import { OrganizationOwnerAutocomplete } from '@spaces/organization/organization-details/owner/OrganizationOwnerAutocomplete';
import { User } from '@spaces/graphql';

export const OwnerTableCell = ({
  organizationId,
  owner,
}: {
  organizationId: string;
  owner?: Pick<User, 'id' | 'firstName' | 'lastName'> | null;
}) => {
  return <OrganizationOwnerAutocomplete id={organizationId} owner={owner} />;
};
