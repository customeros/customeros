import { useMemo, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';

import { produce } from 'immer';
import { useLocalStorage } from 'usehooks-ts';

import { SortingState } from '@ui/presentation/Table';
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

import { useTableState } from '../state';
import { useGetOrganizationsInfiniteQuery } from './useGetOrganizationsInfiniteQuery';

interface UseOrganizationsPageDataProps {
  sorting: SortingState;
}

export const useOrganizationsPageData = ({
  sorting,
}: UseOrganizationsPageDataProps) => {
  const client = getGraphQLClient();
  const searchParams = useSearchParams();
  const { columnFilters } = useTableState();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const [_, setLastActivePosition] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const preset = searchParams?.get('preset');
  const searchTerm = searchParams?.get('search');
  const {
    owner,
    website,
    forecast,
    organization,
    relationship,
    timeToRenewal,
    // lastTouchpoint,
    renewalLikelihood,
  } = columnFilters;

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
            return ['IS_CUSTOMER', [true]];
          }
          if (preset === 'portfolio') {
            const userId = globalCache?.global_Cache?.user.id;

            return ['OWNER_ID', [userId]];
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
        organization?.isActive &&
        (organization?.value || organization?.showEmpty)
      ) {
        draft.AND.push({
          filter: {
            property: 'NAME',
            value: organization?.showEmpty ? '' : organization?.value,
            operation: ComparisonOperator.Contains,
            caseSensitive: false,
          },
        });
      }

      if (website?.isActive && (website?.value || website?.showEmpty)) {
        draft.AND.push({
          filter: {
            property: 'WEBSITE',
            value: website?.showEmpty ? '' : website?.value,
            operation: ComparisonOperator.Contains,
          },
        });
      }

      if (relationship.isActive && relationship?.value) {
        draft.AND.push({
          filter: {
            property: 'IS_CUSTOMER',
            value: columnFilters.relationship.value,
            operation: ComparisonOperator.In,
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
            property: 'RENEWAL_CYCLE_NEXT',
            value: timeToRenewal?.value,
            operation: ComparisonOperator.Lte,
          },
        });
      }

      if (forecast.isActive && forecast.value) {
        draft.AND.push({
          filter: {
            property: 'FORECAST_AMOUNT',
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
            // includeEmpty: owner.showEmpty,
          },
        });
      }
    });
  }, [
    searchParams?.toString(),
    globalCache?.global_Cache?.user.id,
    organization?.isActive,
    organization?.value,
    organization?.showEmpty,
    website?.isActive,
    website.value,
    website?.showEmpty,
    relationship.isActive,
    relationship?.value.length,
    renewalLikelihood?.isActive,
    renewalLikelihood?.value.length,
    timeToRenewal?.isActive,
    timeToRenewal?.value,
    forecast?.isActive,
    forecast?.value[0],
    forecast?.value[1],
    owner?.isActive,
    owner?.value.length,
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

  const totalCount =
    data?.pages?.[0].dashboardView_Organizations?.totalElements;
  const totalAvailable =
    data?.pages?.[0].dashboardView_Organizations?.totalAvailable;

  const flatData = useMemo(
    () =>
      (data?.pages?.flatMap(
        (o) => o.dashboardView_Organizations?.content,
      ) as Organization[]) || [],
    [
      data,
      organization?.isActive,
      organization?.value,
      website?.isActive,
      website.value,
      relationship.isActive,
      relationship?.value.length,
      renewalLikelihood?.isActive,
      renewalLikelihood?.value.length,
      timeToRenewal?.isActive,
      timeToRenewal?.value,
      forecast?.isActive,
      forecast?.value[0],
      forecast?.value[1],
      owner?.isActive,
      owner?.value.length,
    ],
  );

  const allOrganizationIds = flatData.map((o) => o?.id);

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

  return {
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
