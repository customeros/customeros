import { produce } from 'immer';
import { TableViewDefsStore } from '@store/TableViewDefs/TableViewDefs.store';

import { Filter, SortBy, SortingDirection } from '@graphql/types';
import { getServerGraphQLClient } from '@shared/util/getServerGraphQLClient';
import {
  GetInvoicesQuery,
  GetInvoicesDocument,
  GetInvoicesQueryVariables,
} from '@shared/graphql/getInvoices.generated';

import { Preview } from './src/components/Preview';
import { InvoicesTable } from './src/components/InvoicesTable';

export default async function InvoicesPage({
  searchParams,
}: {
  searchParams: { preset?: string; searchTerm?: string };
}) {
  const client = getServerGraphQLClient();
  const { preset, searchTerm } = searchParams;

  let initialData: GetInvoicesQuery | undefined = undefined;

  try {
    const tableViewDefsRes = await TableViewDefsStore.serverSideBootstrap(
      client,
    );
    const tableViewDefs = tableViewDefsRes?.tableViewDefs ?? [];

    const filters = JSON.parse(
      tableViewDefs.find((t) => t.id === preset)?.filters ?? ('{}' as string),
    );

    const whereAndSort = createFilters({ filters, searchTerm });

    initialData = await client.request<
      GetInvoicesQuery,
      GetInvoicesQueryVariables
    >(GetInvoicesDocument, {
      pagination: { limit: 40, page: 0 },
      ...whereAndSort,
    });
  } catch (e) {
    console.error('Failed to fetch initial Invoices data', e);
  }

  return (
    <>
      <InvoicesTable initialData={initialData} />
      <Preview />
    </>
  );
}

function createFilters({
  filters,
  searchTerm,
}: {
  filters: Filter;
  searchTerm?: string;
}) {
  const sort: SortBy = {
    by: 'INVOICE_DUE_DATE',
    direction: SortingDirection.Desc,
    caseSensitive: false,
  };

  const where = (() => {
    return produce<Filter>(filters, (draft) => {
      if (!draft.AND) {
        draft.AND = [];
      }

      if (searchTerm) {
        draft.AND.push({
          filter: {
            property: 'CONTRACT_NAME',
            value: searchTerm,
          },
        });
      }
    });
  })();

  return { sort, where };
}
