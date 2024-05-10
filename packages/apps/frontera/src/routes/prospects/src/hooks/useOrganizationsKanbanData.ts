import { useRef, useMemo, useEffect } from 'react';

import { produce } from 'immer';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { SortingState, TableInstance } from '@ui/presentation/Table';
import { SortBy, Organization, SortingDirection } from '@graphql/types';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';

import { useGetOrganizationsKanbanInfiniteQuery } from './useGetOrganizationsKanbanInfiniteQuery';

interface UseOrganizationsPageDataProps {
  sorting: SortingState;
}

export const useOrganizationsKanbanData = ({
  sorting,
}: UseOrganizationsPageDataProps) => {
  const client = getGraphQLClient();

  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const tableRef = useRef<TableInstance<Organization> | null>(null);

  const sortBy: SortBy | undefined = useMemo(() => {
    if (!sorting.length) return;

    return {
      by: sorting[0].id,
      direction: sorting[0].desc ? SortingDirection.Desc : SortingDirection.Asc,
      caseSensitive: false,
    };
  }, [sorting]);

  const { data, isFetching, isLoading, hasNextPage, fetchNextPage } =
    useGetOrganizationsKanbanInfiniteQuery(client, {
      pagination: {
        page: 1,
        limit: 80,
      },
      sort: sortBy,
    });

  const totalCount =
    data?.pages?.[0].dashboardView_Organizations?.totalElements;
  const totalAvailable =
    data?.pages?.[0].dashboardView_Organizations?.totalAvailable;

  const flatData = useMemo(
    () =>
      (data?.pages?.flatMap(
        (o) => o.dashboardView_Organizations?.content,
      ) as Organization[]) || [],
    [data],
  );

  const allOrganizationIds = flatData.map((o) => o?.metadata.id);

  useEffect(() => {
    setOrganizationsMeta(
      produce(organizationsMeta, (draft) => {
        draft.getOrganization.pagination.page = 1;
        draft.getOrganization.pagination.limit = 80;
        draft.getOrganization.sort = sortBy;
      }),
    );

    tableRef.current?.resetRowSelection();
  }, [sortBy, data?.pageParams]);

  return {
    tableRef,
    isLoading,
    isFetching,
    totalCount,
    hasNextPage,
    fetchNextPage,
    data: flatData,
    totalAvailable,
    allOrganizationIds,
  };
};
