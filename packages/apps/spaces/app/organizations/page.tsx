'use client';

import { useSearchParams } from 'next/navigation';
import { useMemo, useState, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';
import { useIsRestoring } from '@tanstack/react-query';

import { GridItem } from '@ui/layout/Grid';
import { Table, SortingState } from '@ui/presentation/Table';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import {
  Filter,
  SortBy,
  Organization,
  SortingDirection,
  ComparisonOperator,
} from '@graphql/types';

import { useTableState } from './src/state';
import { Search } from './src/components/Search';
import { TableActions } from './src/components/Actions';
import { useOrganizationsPageMethods } from './src/hooks';
import { getColumns } from './src/components/Columns/Columns';
// import EmptyState from './src/components/EmptyState/EmptyState';
import { useGetOrganizationsInfiniteQuery } from './src/hooks/useGetOrganizationsInfiniteQuery';

export default function OrganizationsPage() {
  const client = getGraphQLClient();
  const searchParams = useSearchParams();
  const [tabs, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');
  // const router = useRouter();

  const [sorting, setSorting] = useState<SortingState>([]);

  const { columnFilters } = useTableState();
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const { createOrganization } = useOrganizationsPageMethods();
  const isRestoring = useIsRestoring();

  const { data: globalCache } = useGlobalCacheQuery(client);

  const where = useMemo(() => {
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

      if (preset) {
        const [property, value] = (() => {
          if (preset === 'customer') {
            return ['IS_CUSTOMER', true];
          }
          if (preset === 'portfolio') {
            const userId = globalCache?.global_Cache?.user.id;

            return ['OWNER_ID', userId];
          }

          return [];
        })();
        if (!property || !value) return;
        draft.AND.push({
          filter: {
            property,
            value,
            operation: ComparisonOperator.Eq,
          },
        });
      }

      if (
        columnFilters?.organization?.isActive &&
        columnFilters?.organization?.value
      ) {
        draft.AND.push({
          filter: {
            property: 'ORGANIZATION',
            value: columnFilters?.organization?.value,
            operation: ComparisonOperator.Contains,
            caseSensitive: false,
          },
        });
      }

      if (
        columnFilters?.relationship.isActive &&
        columnFilters?.relationship.value.length > 0 &&
        columnFilters?.relationship.value.length < 2
      ) {
        draft.AND.push({
          filter: {
            property: 'IS_CUSTOMER',
            value: columnFilters.relationship.value[0],
            operation: ComparisonOperator.Eq,
          },
        });
      }
    });
  }, [
    searchParams?.toString(),
    globalCache?.global_Cache?.user.id,
    columnFilters?.organization?.value,
    columnFilters?.relationship.isActive,
    columnFilters?.organization?.isActive,
    columnFilters.relationship.value.length,
  ]);

  const sortBy: SortBy | undefined = useMemo(() => {
    if (!sorting.length) return;

    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  }, [sorting]);
  const { data, isFetching, isLoading, hasNextPage, fetchNextPage } =
    useGetOrganizationsInfiniteQuery(
      client,
      {
        pagination: {
          page: 1,
          limit: 40,
        },
        sort: sortBy,
        where,
      },
      {
        enabled:
          preset === 'portfolio' ? !!globalCache?.global_Cache?.user.id : true,
      },
    );

  const flatData = useMemo(
    () =>
      (data?.pages?.flatMap(
        (o) => o.dashboardView_Organizations?.content,
      ) as Organization[]) || [],
    [
      data,
      columnFilters?.organization?.value,
      columnFilters?.relationship.isActive,
      columnFilters?.organization?.isActive,
      columnFilters.relationship.value.length,
    ],
  );

  const allOrganizationIds = flatData.map((o) => o?.id);

  const handleCreateOrganization = useCallback(() => {
    createOrganization.mutate({ input: { name: '' } });
  }, [createOrganization]);

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const columns = useMemo(
    () =>
      getColumns({
        tabs,
        createIsLoading: createOrganization.isLoading,
        onCreateOrganization: handleCreateOrganization,
      }),
    [tabs, handleCreateOrganization, createOrganization, getColumns],
  );

  useEffect(() => {
    setOrganizationsMeta(
      produce(organizationsMeta, (draft) => {
        draft.getOrganization.pagination.page = 1;
        draft.getOrganization.pagination.limit = 40;
        draft.getOrganization.sort = sortBy;
        draft.getOrganization.where = where;
      }),
    );
    setLastActivePosition((prev) =>
      produce(prev, (draft) => {
        if (!draft?.root) return;
        draft.root = `organizations?${searchParams?.toString()}`;
      }),
    );
  }, [sortBy, searchParams?.toString(), data?.pageParams]);

  // this should only enter if no rows are present for this tenant in db
  // if (
  //   data?.pages?.[0].dashboardView_Organizations?.totalElements === 0 &&
  //   !searchTerm
  // ) {
  //   const emptyOrganization = {
  //     title: "Let's get started",
  //     description:
  //       'Start seeing your customer conversations all in one place by adding an organization',
  //     buttonLabel: 'Add Organization',
  //     onClick: handleCreateOrganization,
  //   };
  //   const myPortfolioEmpty = {
  //     title: 'No organizations assigned to you yet',
  //     description:
  //       'Currently, you have not been assigned to any organizations.\n' +
  //       '\n' +
  //       'Head to your list of organizations and assign yourself as an owner to one of them.',
  //     buttonLabel: 'Go to Organizations',
  //     onClick: () => {
  //       router.push(`/organizations`);
  //     },
  //   };
  //   const emptyState =
  //     preset === 'portfolio' ? myPortfolioEmpty : emptyOrganization;

  //   return <EmptyState {...emptyState} />;
  // }

  return (
    <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
      <Search
        key={preset}
        placeholder={
          preset === 'customer'
            ? 'Search customers'
            : preset === 'portfolio'
            ? 'Search portfolio'
            : 'Search organizations'
        }
      />

      <Table<Organization>
        data={flatData}
        columns={columns}
        sorting={sorting}
        enableTableActions
        enableRowSelection
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        isLoading={isRestoring ? false : isLoading}
        totalItems={
          isRestoring
            ? 40
            : data?.pages?.[0].dashboardView_Organizations?.totalElements || 0
        }
        renderTableActions={(table) => (
          <TableActions
            table={table}
            allOrganizationsIds={allOrganizationIds}
          />
        )}
      />
    </GridItem>
  );
}
