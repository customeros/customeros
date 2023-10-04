'use client';

import { useMemo, useState, useEffect, useCallback } from 'react';
import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';
import { useSearchParams } from 'next/navigation';

import {
  Table,
  SortingState,
  TableInstance,
  RowSelectionState,
} from '@ui/presentation/Table';
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
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

import { Search } from './src/components/Search';
import { useOrganizationsMeta } from './src/shared/state';
import { useOrganizationsPageMethods } from './src/hooks';
import { columns } from './src/components/Columns/Columns';
import EmptyState from './src/components/EmptyState/EmptyState';
import { OrganizationListActions } from './src/components/Actions';
import { useInfiniteGetOrganizationsQuery } from './src/graphql/getOrganizations.generated';

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
  const [enableSelection, setEnableSelection] = useState(false);
  const [selection, setSelection] = useState<RowSelectionState>({});
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const { createOrganization, hideOrganizations, mergeOrganizations } =
    useOrganizationsPageMethods({
      selection,
      setEnableSelection,
    });

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
    setOrganizationsMeta(
      produce(organizationsMeta, (draft) => {
        draft.getOrganization.pagination.page = 1;
      }),
    );

    if (!sorting.length) return;
    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  }, [sorting]);

  const { data, fetchNextPage, isInitialLoading, isFetchingNextPage } =
    useInfiniteGetOrganizationsQuery(
      'pagination',
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
        getNextPageParam: () => {
          return {
            pagination: {
              page: organizationsMeta.getOrganization.pagination.page + 1,
              limit: 40,
            },
          };
        },
      },
    );

  const flatData =
    (data?.pages?.flatMap(
      (o) => o.dashboardView_Organizations?.content,
    ) as Organization[]) || [];
  const allOrganizationIds = flatData.map((o) => o?.id);
  const selectedIds = Object.keys(selection)
    .map((k) => (allOrganizationIds as string[])[Number(k)])
    .filter(Boolean);

  const handleCreateOrganization = () => {
    createOrganization.mutate({ input: { name: '' } });
  };

  const handleFetchMore = useCallback(() => {
    setOrganizationsMeta(
      produce(organizationsMeta, (draft) => {
        draft.getOrganization.pagination.page += 1;
      }),
    );

    fetchNextPage();
  }, [setOrganizationsMeta, fetchNextPage, organizationsMeta]);

  const handleMergeOrganizations = (_: TableInstance<Organization>) => {
    const primaryId = selectedIds[0];
    const mergeIds = selectedIds.slice(1);

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
    setEnableSelection(false);
  };

  const handleCancelRemoveOrganizations = () => {
    onClose();
    setEnableSelection(false);
  };

  useEffect(() => {
    setOrganizationsMeta(
      produce(organizationsMeta, (draft) => {
        draft.getOrganization.pagination.page = 1;
        draft.getOrganization.pagination.limit = 40;
        draft.getOrganization.sort = sortBy;
        draft.getOrganization.where = where;
      }),
    );
  }, [sortBy, where]);

  useEffect(() => {
    setLastActivePosition((prev) =>
      produce(prev, (draft) => {
        if (!draft?.root) return;
        draft.root = `organizations?${searchParams?.toString()}`;
      }),
    );
  }, [searchParams?.toString()]);

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
        columns={columns(tabs)}
        sorting={sorting}
        enableTableActions
        isLoading={isInitialLoading || isFetchingNextPage}
        selection={selection}
        onSortingChange={setSorting}
        onFetchMore={handleFetchMore}
        totalItems={
          data?.pages?.[0].dashboardView_Organizations?.totalElements || 0
        }
        onSelectionChange={setSelection}
        enableRowSelection={enableSelection}
        renderTableActions={(table) => (
          <OrganizationListActions
            table={table}
            selection={selection}
            onArchiveOrganizations={onOpen}
            isSelectionEnabled={enableSelection}
            toggleSelection={setEnableSelection}
            onCreateOrganization={handleCreateOrganization}
            onMergeOrganizations={handleMergeOrganizations}
          />
        )}
      />

      <ConfirmDeleteDialog
        isOpen={isOpen}
        icon={<Archive />}
        confirmButtonLabel={'Archive'}
        onConfirm={handleHideOrganizations}
        onClose={handleCancelRemoveOrganizations}
        isLoading={hideOrganizations.isLoading}
        label={`Archive selected ${
          selectedIds.length === 1 ? 'organization' : 'organizations'
        }?`}
      />
    </GridItem>
  );
}
