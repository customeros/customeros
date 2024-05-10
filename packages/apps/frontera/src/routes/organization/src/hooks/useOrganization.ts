import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationQuery } from '@organization/graphql/organization.generated';

export const useOrganization = ({ id }: { id: string }) => {
  const client = getGraphQLClient();

  return useOrganizationQuery(client, { id });
};
