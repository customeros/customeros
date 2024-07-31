import { useSearchParams } from 'react-router-dom';
import { useRef, useMemo, useEffect } from 'react';

import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useRenewalsMeta } from '@shared/state/RenewalsMeta.atom';
import { SortingState, TableInstance } from '@ui/presentation/Table';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import {
  Filter,
  SortBy,
  RenewalRecord,
  SortingDirection,
  ComparisonOperator,
} from '@graphql/types';

import { useTableState } from '../state';
import { useGetRenewalsInfiniteQuery } from './useGetRenewalsInfiniteQuery';

interface UseRenewalsPageDataProps {
  sorting: SortingState;
}

export const useRenewalsPageData = ({ sorting }: UseRenewalsPageDataProps) => {
  const client = getGraphQLClient();
  const [searchParams] = useSearchParams();
  const { columnFilters } = useTableState();
  const store = useStore();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const [renewalsMeta, setRenewalsMeta] = useRenewalsMeta();
  const [_, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'renewals' });
  const tableRef = useRef<TableInstance<RenewalRecord> | null>(null);

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');
  const {
    owner,
    forecast,
    organization,
    timeToRenewal,
    lastTouchpoint,
    renewalLikelihood,
  } = columnFilters;

  const where = useMemo(() => {
    const defaultFilters = JSON.parse(
      store.tableViewDefs.getById(preset ?? '1')?.value.filters ?? '{}',
    );

    return produce<Filter>(defaultFilters, (draft) => {
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

      if (
        organization?.isActive &&
        (organization?.value || organization?.showEmpty)
      ) {
        draft.AND.push({
          filter: {
            property: 'NAME',
            value: organization?.value,
            operation: ComparisonOperator.Contains,
            caseSensitive: false,
            includeEmpty: organization?.showEmpty,
          },
        });
      }

      if (renewalLikelihood?.isActive && renewalLikelihood?.value) {
        draft.AND.push({
          filter: {
            property: 'RENEWAL_LIKELIHOOD',
            value: renewalLikelihood?.value,
            operation: ComparisonOperator.In,
          },
        });
      }

      if (timeToRenewal?.isActive && timeToRenewal?.value) {
        draft.AND.push({
          filter: {
            property: 'RENEWAL_DATE',
            value: timeToRenewal?.value,
            operation: ComparisonOperator.Lte,
          },
        });
      }

      if (forecast.isActive && forecast.value) {
        draft.AND.push({
          filter: {
            property: 'FORECAST_ARR',
            value: forecast.value,
            operation: ComparisonOperator.Between,
          },
        });
      }

      if (owner.isActive && owner.value) {
        draft.AND.push({
          filter: {
            property: 'OWNER_ID',
            value: owner.value,
            operation: ComparisonOperator.In,
            includeEmpty: owner.showEmpty,
          },
        });
      }

      if (lastTouchpoint.isActive) {
        if (lastTouchpoint.value.length) {
          draft.AND.push({
            filter: {
              property: 'LAST_TOUCHPOINT_TYPE',
              value: lastTouchpoint.value,
              operation: ComparisonOperator.In,
            },
          });
        }

        if (lastTouchpoint.after) {
          draft.AND.push({
            filter: {
              property: 'LAST_TOUCHPOINT_AT',
              value: lastTouchpoint.after,
              operation: ComparisonOperator.Gte,
            },
          });
        }
      }
    });
  }, [
    searchParams?.toString(),
    globalCache?.global_Cache?.user.id,
    organization?.isActive,
    organization?.value,
    organization?.showEmpty,
    renewalLikelihood?.isActive,
    renewalLikelihood?.value.length,
    timeToRenewal?.isActive,
    timeToRenewal?.value,
    forecast?.isActive,
    forecast?.value[0],
    forecast?.value[1],
    owner?.isActive,
    owner?.value.length,
    owner?.showEmpty,
    lastTouchpoint?.isActive,
    lastTouchpoint?.value,
    lastTouchpoint?.after,
  ]);

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
    isLoading,
    hasNextPage,
    fetchNextPage,
    isRefetching,
  } = useGetRenewalsInfiniteQuery(client, {
    pagination: {
      page: 1,
      limit: 40,
    },
    sort: sortBy,
    where,
  });

  const totalCount = data?.pages?.[0].dashboardView_Renewals?.totalElements;
  const totalAvailable =
    data?.pages?.[0].dashboardView_Renewals?.totalAvailable;

  const flatData = useMemo(
    () =>
      (data?.pages?.flatMap(
        (o) => o.dashboardView_Renewals?.content,
      ) as RenewalRecord[]) || [],
    [
      data,
      organization?.isActive,
      organization?.value,
      organization?.showEmpty,
      renewalLikelihood?.isActive,
      renewalLikelihood?.value.length,
      timeToRenewal?.isActive,
      timeToRenewal?.value,
      forecast?.isActive,
      forecast?.value[0],
      forecast?.value[1],
      owner?.isActive,
      owner?.value.length,
      owner?.showEmpty,
      lastTouchpoint?.isActive,
      lastTouchpoint?.value.length,
      lastTouchpoint?.after,
    ],
  );

  useEffect(() => {
    setRenewalsMeta(
      produce(renewalsMeta, (draft) => {
        draft.getRenewals.pagination.page = 1;
        draft.getRenewals.pagination.limit = 40;
        draft.getRenewals.sort = sortBy;
        draft.getRenewals.where = where;
      }),
    );
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
    isRefetching,
    fetchNextPage,
    data: flatData,
    totalAvailable,
  };
};
