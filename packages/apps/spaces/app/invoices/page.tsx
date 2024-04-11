import { produce } from 'immer';

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
    const filters = createFilters({ preset, searchTerm });

    initialData = await client.request<
      GetInvoicesQuery,
      GetInvoicesQueryVariables
    >(GetInvoicesDocument, { pagination: { limit: 40, page: 0 }, ...filters });
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
  preset,
  searchTerm,
}: {
  preset?: string;
  searchTerm?: string;
}) {
  const sort: SortBy = {
    by: 'INVOICE_DUE_DATE',
    direction: SortingDirection.Desc,
    caseSensitive: false,
  };

  const where = (() => {
    if (!preset) return undefined;

    return produce<Filter>({ AND: [] }, (draft) => {
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

      if (preset) {
        switch (preset) {
          case '4':
            draft.AND.push({
              filter: {
                property: 'INVOICE_PREVIEW',
                value: true,
              },
            });
            break;
          case '5':
            draft.AND.push({
              filter: {
                property: 'INVOICE_DRY_RUN',
                value: false,
              },
            });
            break;
          default:
            break;
        }
      }
    });
  })();

  return { sort, where };
}
