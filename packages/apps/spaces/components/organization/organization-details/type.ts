import { GetOrganizationQuery } from '@spaces/graphql';

export interface OrganizationDetailsProps {
  id: string;
  loading: boolean;
  organization: GetOrganizationQuery['organization'] | undefined | null;
}
