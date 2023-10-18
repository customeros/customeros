'use client';

import { useMemo, useState, useEffect, useCallback } from 'react';
import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';
import { useSearchParams } from 'next/navigation';
import { useIsRestoring } from '@tanstack/react-query';

import {
  Filter,
  SortBy,
  Organization,
  SortingDirection,
  ComparisonOperator,
} from '@graphql/types';
import { GridItem } from '@ui/layout/Grid';
import { Table, SortingState } from '@ui/presentation/Table';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { Search } from './src/components/Search';
import { TableActions } from './src/components/Actions';
import { useOrganizationsPageMethods } from './src/hooks';
import { getColumns } from './src/components/Columns/Columns';
import EmptyState from './src/components/EmptyState/EmptyState';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useGetOrganizationsInfiniteQuery } from './src/hooks/useGetOrganizationsInfiniteQuery';

export default function OrganizationsPage() {
  const client = getGraphQLClient();
  const searchParams = useSearchParams();
  const [tabs, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');

  const [sorting, setSorting] = useState<SortingState>([]);
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const { createOrganization } = useOrganizationsPageMethods();
  const isRestoring = useIsRestoring();

  const { data: globalCache } = useGlobalCacheQuery(client);

  const where = useMemo(() => {
    return produce<Filter>({ AND: [] }, (draft) => {
      if (!draft.AND) draft.AND = [];
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
            const userId = globalCache?.global_Cache.user.id;
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
    });
  }, [searchParams?.toString()]);

  const sortBy: SortBy | undefined = useMemo(() => {
    if (!sorting.length) return;
    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  }, [sorting]);

  const { data, isLoading, isFetching, hasNextPage, fetchNextPage } =
    useGetOrganizationsInfiniteQuery(client, {
      pagination: {
        page: 1,
        limit: 40,
      },
      sort: sortBy,
      where,
    });

  const flatData = useMemo(
    () =>
      (data?.pages?.flatMap(
        (o) => o.dashboardView_Organizations?.content,
      ) as Organization[]) || [],
    [data],
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

  if (
    data?.pages?.[0].dashboardView_Organizations?.totalElements === 0 &&
    !searchTerm
  ) {
    return <EmptyState onClick={handleCreateOrganization} />;
  }

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
        isLoading={isRestoring ? false : isLoading}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
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
