import { useMemo, useEffect } from 'react';

import { produce } from 'immer';
import {
  InfiniteData,
  FetchNextPageOptions,
  InfiniteQueryObserverResult,
} from '@tanstack/react-query';

import { Invoice } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { filterOutDryRunInvoices } from '@shared/components/Invoice/utils';
import { useOrganizationInvoicesMeta } from '@shared/state/OrganizationInvoicesMeta.atom';
import {
  GetInvoicesQuery,
  useInfiniteGetInvoicesQuery,
} from '@shared/graphql/getInvoices.generated';

export type InvoiceTableData = Invoice & { id: string };
interface useInfiniteInvoicesReturn {
  isFetched: boolean;
  isFetching: boolean;
  hasNextPage: boolean;
  totalInvoicesCount: number;
  invoiceFlattenData: InvoiceTableData[];
  fetchNextPage: (
    options?: FetchNextPageOptions | undefined,
  ) => Promise<
    InfiniteQueryObserverResult<
      InfiniteData<GetInvoicesQuery, unknown>,
      unknown
    >
  >;
}

export function useInfiniteInvoices(
  organizationId?: string,
): useInfiniteInvoicesReturn {
  const client = getGraphQLClient();
  const [invoicesMeta, setInvoicesMeta] = useOrganizationInvoicesMeta();

  const { data, isFetching, isFetched, fetchNextPage, hasNextPage } =
    useInfiniteGetInvoicesQuery(
      client,
      {
        pagination: { page: 0, limit: 40 },
        where: { ...filterOutDryRunInvoices },
        organizationId,
      },
      {
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
            pagination: { page: allPages.length + 1, limit: 40 },
            organizationId,
            where: { ...filterOutDryRunInvoices },
          };
        },
      },
    );

  const invoiceFlattenData = useMemo(() => {
    return (
      (data?.pages?.flatMap((o) => o.invoices?.content) as Invoice[])?.map(
        (e) => ({ ...e, id: e.metadata.id }),
      ) || []
    );
  }, [data]);

  const totalInvoicesCount = useMemo(() => {
    return data?.pages?.[0]?.invoices?.totalElements ?? 0;
  }, [data?.pages?.[0]?.invoices?.totalElements]);

  useEffect(() => {
    setInvoicesMeta(
      produce(invoicesMeta, (draft) => {
        draft.getInvoices.pagination.page = 0;
        draft.getInvoices.pagination.limit = 40;

        if (organizationId) {
          draft.getInvoices.organizationId = organizationId;
        }
      }),
    );
  }, [organizationId, data?.pageParams]);

  return {
    invoiceFlattenData,
    totalInvoicesCount,
    isFetching,
    isFetched,
    fetchNextPage,
    hasNextPage,
  };
}
