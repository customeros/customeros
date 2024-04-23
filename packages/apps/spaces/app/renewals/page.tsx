import { ErrorBoundary } from 'react-error-boundary';

import { TableViewDefsStore } from '@store/TableViewDefs/TableViewDefs.store';

import { getServerGraphQLClient } from '@shared/util/getServerGraphQLClient';

import { RenewalsTable } from './src/components/RenewalsTable';
import {
  GetRenewalsDocument,
  type GetRenewalsQuery,
  type GetRenewalsQueryVariables,
} from './src/graphql/getRenewals.generated';

export default async function RenewalsPage({
  searchParams,
}: {
  searchParams?: { [key: string]: string | string[] | undefined };
}) {
  const client = getServerGraphQLClient();
  const preset = searchParams?.preset ?? '1';

  try {
    const tableViewDefsRes = await TableViewDefsStore.serverSideBootstrap(
      client,
    );
    const tableViewDefs = tableViewDefsRes?.tableViewDefs ?? [];

    const where = JSON.parse(
      tableViewDefs.find((t) => t.id === preset)?.filters ?? ('{}' as string),
    );

    // this should be replaced with serverSideBootstrap call when store's implemented.
    const renewals = await client.request<
      GetRenewalsQuery,
      GetRenewalsQueryVariables
    >(GetRenewalsDocument, { pagination: { page: 1, limit: 40 }, where });

    return (
      <RenewalsTable
        bootstrap={{
          renewals,
          tableViewDefs,
        }}
      />
    );
  } catch (e) {
    return <ErrorBoundary fallback={<div>Error hapened</div>} />;
  }
}
