import { produce } from 'immer';

import { Filter } from '@graphql/types';
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

  const res = await client.request<GetInvoicesQuery, GetInvoicesQueryVariables>(
    GetInvoicesDocument,
    { pagination: { limit: 40, page: 0 }, where },
  );

  return (
    <>
      <InvoicesTable initialData={res} />
      <Preview />
    </>
  );
}
