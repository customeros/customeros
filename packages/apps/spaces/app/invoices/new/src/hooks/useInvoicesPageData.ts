import { useSearchParams } from 'next/navigation';
import { useRef, useMemo, useEffect } from 'react';

import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
// import { useRenewalsMeta } from '@shared/state/RenewalsMeta.atom';
import { SortingState, TableInstance } from '@ui/presentation/Table';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import {
  Filter,
  SortBy,
  Invoice,
  SortingDirection,
  // ComparisonOperator,
} from '@graphql/types';

// import { useTableState } from '../state';
import { useGetInvoicesInfiniteQuery } from './useGetInvoicesInfiniteQuery';

interface UseRenewalsPageDataProps {
  sorting: SortingState;
}

export const useInvoicesPageData = ({ sorting }: UseRenewalsPageDataProps) => {
  const client = getGraphQLClient();
  const searchParams = useSearchParams();
  // const { columnFilters } = useTableState();
  const { data: globalCache } = useGlobalCacheQuery(client);
  // const [renewalsMeta, setRenewalsMeta] = useRenewalsMeta();
  const [_, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'invoices' });
  const tableRef = useRef<TableInstance<Invoice> | null>(null);

  const preset = searchParams?.get('preset');
  // const searchTerm = searchParams?.get('search');
  // const {
  //   owner,
  //   forecast,
  //   organization,
  //   timeToRenewal,
  //   lastTouchpoint,
  //   renewalLikelihood,
  // } = columnFilters;

  const where = useMemo(() => {
    if (preset === '5') {
      return {
        filter: {
          property: 'PREVIEW',
          value: true,
        },
      } as Filter;
    }

    return undefined;
  }, [searchParams?.toString(), globalCache?.global_Cache?.user.id]);

  const sortBy: SortBy | undefined = useMemo(() => {
    if (!sorting.length) return;

    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  }, [sorting]);

  const { data, isFetching, isLoading, hasNextPage, fetchNextPage } =
    useGetInvoicesInfiniteQuery(client, {
      pagination: {
        page: 0,
        limit: 40,
      },
      // sort: sortBy,
      where,
    });

  const totalCount = data?.pages?.[0].invoices?.totalElements;
  // const totalAvailable = data?.pages?.[0].invoices?.totalAvailable;
  const totalAvailable = (data?.pages?.[0].invoices?.totalPages ?? 0) * 40;

  const flatData = useMemo(
    () => (data?.pages?.flatMap((o) => o.invoices?.content) as Invoice[]) || [],
    [
      data,
      // organization?.isActive,
      // organization?.value,
      // organization?.showEmpty,
      // renewalLikelihood?.isActive,
      // renewalLikelihood?.value.length,
      // timeToRenewal?.isActive,
      // timeToRenewal?.value,
      // forecast?.isActive,
      // forecast?.value[0],
      // forecast?.value[1],
      // owner?.isActive,
      // owner?.value.length,
      // owner?.showEmpty,
      // lastTouchpoint?.isActive,
      // lastTouchpoint?.value.length,
      // lastTouchpoint?.after,
    ],
  );

  useEffect(() => {
    // setRenewalsMeta(
    //   produce(renewalsMeta, (draft) => {
    //     draft.getRenewals.pagination.page = 1;
    //     draft.getRenewals.pagination.limit = 40;
    //     draft.getRenewals.sort = sortBy;
    //     // draft.getRenewals.where = where;
    //   }),
    // );
    setLastActivePosition((prev) =>
      produce(prev, (draft) => {
        if (!draft?.root) return;
        draft.root = `renewals?${searchParams?.toString()}`;
      }),
    );

    tableRef.current?.resetRowSelection();
  }, [sortBy, searchParams?.toString(), data?.pageParams]);

  return {
    tableRef,
    isLoading,
    isFetching,
    totalCount,
    hasNextPage,
    fetchNextPage,
    data: flatData,
    totalAvailable,
  };
};
