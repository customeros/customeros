import { HydrationBoundary } from '@tanstack/react-query';

import { getDehydratedState } from '@shared/util/getDehydratedState';
import {
  useGetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

import { KMenu } from './src/components/KMenu';
import { Search } from './src/components/Search';
import { OrganizationsTable } from './src/components/OrganizationsTable';

export default async function OrganizationsPage() {
  const dehydratedClient = await getDehydratedState(
    useInfiniteGetOrganizationsQuery,
    {
      variables: { page: 0, limit: 40 },
      fetcher: useGetOrganizationsQuery.fetcher,
    },
  );

  return (
    <HydrationBoundary state={dehydratedClient}>
      <Search />
      <OrganizationsTable />
      <KMenu />
    </HydrationBoundary>
  );
}
