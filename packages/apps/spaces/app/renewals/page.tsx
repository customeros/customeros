import { match } from 'ts-pattern';
import { HydrationBoundary } from '@tanstack/react-query';

import { getDehydratedState } from '@shared/util/getDehydratedState';

import { Search } from './src/components/Search';
import { RenewalsTable } from './src/components/RenewalsTable';
import {
  useGetRenewalsQuery,
  useInfiniteGetRenewalsQuery,
} from './src/graphql/getRenewals.generated';

export default async function RenewalsPage({
  searchParams,
}: {
  searchParams?: { [key: string]: string | string[] | undefined };
}) {
  const preset = searchParams?.preset ?? '1';
  const value = match(preset)
    .with('1', () => 'MONTHLY')
    .with('2', () => 'QUARTERLY')
    .with('3', () => 'ANNUALLY')
    .otherwise(() => 'MONTHLY');

  const where = {
    AND: [
      {
        filter: {
          property: 'RENEWAL_CYCLE',
          value,
          operation: 'EQ',
          includeEmpty: false,
        },
      },
    ],
  };

  const dehydratedClient = await getDehydratedState(
    useInfiniteGetRenewalsQuery,
    {
      variables: { pagination: { page: 0, limit: 40 }, where },
      fetcher: useGetRenewalsQuery.fetcher,
    },
  );

  return (
    <HydrationBoundary state={dehydratedClient}>
      <Search />
      <RenewalsTable />
    </HydrationBoundary>
  );
}
