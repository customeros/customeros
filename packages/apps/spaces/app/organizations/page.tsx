import { produce } from 'immer';

import { Filter, ComparisonOperator } from '@graphql/types';
import { getServerGraphQLClient } from '@shared/util/getServerGraphQLClient';
import {
  GetOrganizationsQuery,
  GetOrganizationsDocument,
  GetOrganizationsQueryVariables,
} from '@organizations/graphql/getOrganizations.generated';

import { KMenu } from './src/components/KMenu';
import { Search } from './src/components/Search';
import { OrganizationsTable } from './src/components/OrganizationsTable';

export default async function OrganizationsPage({
  searchParams,
}: {
  searchParams: { searchTerm?: string };
}) {
  const client = getServerGraphQLClient();
  const { searchTerm } = searchParams;

  const where = (() => {
    return produce<Filter>({ AND: [] }, (draft) => {
      if (!draft.AND) {
        draft.AND = [];
      }

      if (searchTerm) {
        draft.AND.push({
          filter: {
            property: 'ORGANIZATION',
            value: searchTerm,
            operation: ComparisonOperator.Contains,
            caseSensitive: false,
          },
        });
      }
    });
  })();

  let initialData: GetOrganizationsQuery | undefined = undefined;

  try {
    initialData = await client.request<
      GetOrganizationsQuery,
      GetOrganizationsQueryVariables
    >(GetOrganizationsDocument, {
      pagination: { limit: 40, page: 1 },
      where,
    });
  } catch (e) {
    console.error('Failed to fetch initial Organizations data', e);
  }

  return (
    <>
      <Search />
      <OrganizationsTable initialData={initialData} />
      <KMenu />
    </>
  );
}
