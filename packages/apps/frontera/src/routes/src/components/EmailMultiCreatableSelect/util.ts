import { QueryClient } from '@tanstack/react-query';

import { useOrganizationPeoplePanelQuery } from '@organization/graphql/organizationPeoplePanel.generated';

export function invalidateQuery(queryClient: QueryClient, id: string) {
  queryClient.invalidateQueries({
    queryKey: useOrganizationPeoplePanelQuery.getKey({ id }),
  });
}
