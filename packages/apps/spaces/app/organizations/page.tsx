'use client';

import { useMemo, useState, useEffect, useCallback } from 'react';
import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';
import { useSearchParams } from 'next/navigation';

import {
  Filter,
  SortBy,
  Organization,
  SortingDirection,
  ComparisonOperator,
} from '@graphql/types';
import { useDisclosure } from '@ui/utils';
import { GridItem } from '@ui/layout/Grid';
import { Archive } from '@ui/media/icons/Archive';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { Table, SortingState, RowSelectionState } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

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
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [tabs, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');

  const [sorting, setSorting] = useState<SortingState>([]);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [targetSelection, setTargetSelection] = useState<
    [index: number, id: string] | null
  >(null);
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const { createOrganization, hideOrganizations, mergeOrganizations } =
    useOrganizationsPageMethods({ selection, setSelection, targetSelection });

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
            return ['RELATIONSHIP', preset.toUpperCase()];
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

  const {
    data,
    isFetching,
    hasNextPage,
    fetchNextPage,
    isInitialLoading,
    isFetchingNextPage,
  } = useGetOrganizationsInfiniteQuery(client, {
    pagination: {
      page: 1,
      limit: 40,
    },
    sort: sortBy,
    where,
  });

  const flatData =
    (data?.pages?.flatMap(
      (o) => o.dashboardView_Organizations?.content,
    ) as Organization[]) || [];
  const allOrganizationIds = flatData.map((o) => o?.id);
  const selectedIds = Object.keys(selection).map(
    (k) => (allOrganizationIds as string[])[Number(k)],
  );

  const handleCreateOrganization = () => {
    createOrganization.mutate({ input: { name: '' } });
  };

  const handleFetchMore = useCallback(() => {
    !isFetching && fetchNextPage();
  }, [fetchNextPage, isFetching]);

  const handleMergeOrganizations = () => {
    const primaryId = targetSelection?.[1];
    const mergeIds = selectedIds.filter((id) => id !== primaryId);

    if (!primaryId || !mergeIds.length) return;

    mergeOrganizations.mutate({
      primaryOrganizationId: primaryId,
      mergedOrganizationIds: mergeIds,
    });
  };

  const handleHideOrganizations = () => {
    const selectedIds = Object.keys(selection)
      .map((k) => (allOrganizationIds as string[])[Number(k)])
      .filter(Boolean);

    hideOrganizations.mutate({
      ids: selectedIds,
    });
    onClose();
  };

  const columns = useMemo(
    () =>
      getColumns({
        tabs,
        createIsLoading: createOrganization.isLoading,
        onCreateOrganization: handleCreateOrganization,
      }),
    [tabs, handleCreateOrganization, createOrganization.isLoading],
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
  }, [sortBy, searchParams?.toString(), data?.pageParams]);

  useEffect(() => {
    setLastActivePosition((prev) =>
      produce(prev, (draft) => {
        if (!draft?.root) return;
        draft.root = `organizations?${searchParams?.toString()}`;
      }),
    );
  }, [searchParams?.toString()]);

  useEffect(() => {
    if (selectedIds.length === 1) {
      const id = selectedIds[0];
      const index = Number(Object.keys(selection)[0]);
      setTargetSelection([index, id]);
    }
  }, [selectedIds.length]);

  if (
    data?.pages?.[0].dashboardView_Organizations?.totalElements === 0 &&
    !searchTerm
  ) {
    return <EmptyState onClick={handleCreateOrganization} />;
  }

  return (
    <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
      <Search
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
        isLoading={isInitialLoading || isFetchingNextPage}
        selection={selection}
        canFetchMore={hasNextPage}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        onSelectionChange={setSelection}
        totalItems={
          data?.pages?.[0].dashboardView_Organizations?.totalElements || 0
        }
        renderTableActions={(table) => (
          <TableActions
            table={table}
            selection={selection}
            onArchiveOrganizations={onOpen}
            onMergeOrganizations={handleMergeOrganizations}
          />
        )}
      />

      <ConfirmDeleteDialog
        isOpen={isOpen}
        icon={<Archive />}
        onClose={onClose}
        confirmButtonLabel={'Archive'}
        onConfirm={handleHideOrganizations}
        isLoading={hideOrganizations.isLoading}
        label={`Archive selected ${
          selectedIds.length === 1 ? 'organization' : 'organizations'
        }?`}
      />
    </GridItem>
  );
}
