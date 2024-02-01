import { useMemo } from 'react';

import { Invoice } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { filterOutDryRunInvoices } from '@shared/components/Invoice/utils';
import { useInfiniteGetInvoicesQuery } from '@shared/graphql/getInvoices.generated';

interface useInfiniteInvoicesReturn {
  isFetched: boolean;
  isFetching: boolean;
  fetchNextPage: () => void; // Replace with the actual function type if different
  totalInvoicesCount: number;
  invoiceFlattenData: Invoice[];
}
export function useInfiniteInvoices(
  organizationId?: string,
): useInfiniteInvoicesReturn {
  const client = getGraphQLClient();

  const { data, isFetching, isFetched, fetchNextPage } =
    useInfiniteGetInvoicesQuery(
      client,
      {
        pagination: { page: 1, limit: 40 },
        where: { ...filterOutDryRunInvoices },
        organizationId,
      },
      {
        enabled: true,
        initialPageParam: 1,
        getNextPageParam: (lastPage, allPages) => {
          const content = allPages.flatMap(
            (page) => page.invoices?.content ?? [],
          );
          const totalElements = lastPage.invoices?.totalElements ?? 0;

          if (content.length >= totalElements) {
            return undefined;
          }

          return {
            pagination: { page: allPages.length + 1, limit: 5 },
            organizationId,
            where: { ...filterOutDryRunInvoices },
          };
        },
      },
    );

  const invoiceFlattenData = useMemo(
    () => (data?.pages?.flatMap((o) => o.invoices?.content) as Invoice[]) || [],
    [data],
  );

  const totalInvoicesCount = useMemo(
    () => data?.pages?.[0]?.invoices?.totalElements ?? 0,
    [data?.pages?.[0]?.invoices?.totalElements],
  );

  return {
    invoiceFlattenData,
    totalInvoicesCount,
    isFetching,
    isFetched,
    fetchNextPage,
  };
}
