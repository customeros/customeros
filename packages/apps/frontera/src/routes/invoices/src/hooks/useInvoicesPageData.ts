import { useRef, useMemo, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useInvoicesMeta } from '@shared/state/InvoicesMeta.atom';
import { SortingState, TableInstance } from '@ui/presentation/Table';
import { GetInvoicesQuery } from '@shared/graphql/getInvoices.generated';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import {
  Filter,
  SortBy,
  Invoice,
  SortingDirection,
  ComparisonOperator,
} from '@graphql/types';

import { useTableState } from '../state';
import { useGetInvoicesInfiniteQuery } from './useGetInvoicesInfiniteQuery';

interface UseRenewalsPageDataProps {
  sorting: SortingState;
  initialData?: GetInvoicesQuery;
}

export const useInvoicesPageData = ({
  sorting,
  initialData,
}: UseRenewalsPageDataProps) => {
  const client = getGraphQLClient();
  const [searchParams] = useSearchParams();
  const { columnFilters } = useTableState();
  const store = useStore();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const [invoicesMeta, setInvoicesMeta] = useInvoicesMeta();
  const tableRef = useRef<TableInstance<Invoice> | null>(null);

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');
  const { issueDate, billingCycle, invoiceStatus, paymentStatus } =
    columnFilters;

  const [_, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, {
    root: `invoices?preset=${preset}`,
  });

  const filtersStr =
    store.tableViewDefs.getById(preset ?? '1')?.value.filters ?? '{}';

  const where = useMemo(() => {
    const defaultFilters = JSON.parse(filtersStr);

    return produce<Filter>(defaultFilters, (draft) => {
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

      if (billingCycle?.isActive) {
        draft.AND.push({
          filter: {
            property: 'CONTRACT_BILLING_CYCLE',
            value: billingCycle.value,
            operation: ComparisonOperator.In,
          },
        });
      }

      if (issueDate?.isActive) {
        draft.AND.push({
          filter: {
            property: 'INVOICE_ISSUED_DATE',
            value:
              preset === '5'
                ? [issueDate.value, new Date().toISOString()]
                : [new Date().toISOString(), issueDate.value],
            operation: ComparisonOperator.In,
          },
        });
      }

      if (
        invoiceStatus?.isActive &&
        typeof invoiceStatus?.value !== 'undefined'
      ) {
        draft.AND.push({
          filter: {
            property: 'CONTRACT_ENDED',
            value: invoiceStatus.value,
            operation: ComparisonOperator.In,
          },
        });
      }

      if (paymentStatus?.isActive && paymentStatus?.value?.length) {
        draft.AND.push({
          filter: {
            property: 'INVOICE_STATUS',
            value: paymentStatus.value,
            operation: ComparisonOperator.In,
          },
        });
      }
    });
  }, [
    filtersStr,
    searchParams?.toString(),
    globalCache?.global_Cache?.user.id,
    billingCycle?.isActive,
    billingCycle?.value?.length,
    issueDate?.isActive,
    issueDate?.value,
    invoiceStatus?.isActive,
    invoiceStatus?.value,
    paymentStatus?.isActive,
    paymentStatus?.value?.length,
  ]);

  const sortBy: SortBy[] | undefined = useMemo(() => {
    if (!sorting.length) return;

    return [
      {
        by: sorting[0].id,
        direction: sorting[0].desc
          ? SortingDirection.Desc
          : SortingDirection.Asc,
        caseSensitive: false,
      },
    ];
  }, [sorting]);

  const {
    data,
    isFetching,
    isLoading,
    isPending,
    hasNextPage,
    isRefetching,
    fetchNextPage,
  } = useGetInvoicesInfiniteQuery(
    client,
    {
      pagination: {
        page: 0,
        limit: 40,
      },
      sort: sortBy,
      where,
    },
    {
      placeholderData: initialData
        ? { pageParams: [0], pages: [initialData] }
        : undefined,
    },
  );

  const totalCount = data?.pages?.[0].invoices?.totalElements;
  const totalAvailable = data?.pages?.[0].invoices?.totalAvailable;

  const flatData = useMemo(
    () => (data?.pages?.flatMap((o) => o.invoices?.content) as Invoice[]) || [],
    [
      data,
      billingCycle?.isActive,
      billingCycle?.value,
      issueDate?.isActive,
      issueDate?.value,
      invoiceStatus?.isActive,
      invoiceStatus?.value,
      paymentStatus?.isActive,
      paymentStatus?.value,
    ],
  );

  useEffect(() => {
    setInvoicesMeta(
      produce(invoicesMeta, (draft) => {
        draft.getInvoices.pagination.page = 0;
        draft.getInvoices.pagination.limit = 40;
        draft.getInvoices.sort = sortBy;
        draft.getInvoices.where = where;
      }),
    );
    setLastActivePosition((prev) =>
      produce(prev, (draft) => {
        if (!draft?.root) return;
        draft.root = `invoices?${searchParams?.toString()}`;
      }),
    );

    tableRef.current?.resetRowSelection();
  }, [sortBy, searchParams?.toString(), data?.pageParams]);

  return {
    tableRef,
    isLoading,
    isFetching,
    isPending,
    totalCount,
    hasNextPage,
    isRefetching,
    fetchNextPage,
    data: flatData,
    totalAvailable,
  };
};
