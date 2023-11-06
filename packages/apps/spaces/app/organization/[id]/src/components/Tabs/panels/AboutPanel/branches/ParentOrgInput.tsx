import React, { useState } from 'react';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Select } from '@ui/form/SyncSelect';
import { ComparisonOperator } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';
import { useGetOrganizationOptionsQuery } from '@organization/src/graphql/getOrganizationOptions.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/src/graphql/organization.generated';
import { useAddSubsidiaryToOrganizationMutation } from '@organization/src/graphql/addSubsidiaryToOrganization.generated';
import { useRemoveSubsidiaryToOrganizationMutation } from '@organization/src/graphql/removeSubsidiaryToOrganization.generated';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

interface ParentOrgInputProps {
  id: string;
  parentOrg: { label: string; value: string } | null;
}

export const ParentOrgInput: React.FC<ParentOrgInputProps> = ({
  parentOrg,
  id,
}) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState('');
  const [organizationsMeta] = useOrganizationsMeta();
  const { data, isLoading } = useGetOrganizationOptionsQuery(client, {
    pagination: {
      page: 1,
      limit: 30,
    },
    sort: undefined,
    where: {
      filter: {
        property: 'ORGANIZATION',
        value: searchTerm,
        operation: ComparisonOperator.Contains,
        caseSensitive: false,
      },
    },
  });
  const queryKey = useOrganizationQuery.getKey({ id });
  const organizationsQueryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );
  const invalidateQuery = () => queryClient.invalidateQueries(queryKey);

  const addSubsidiaryToOrganizationMutation =
    useAddSubsidiaryToOrganizationMutation(client, {
      onMutate: ({ input }) => {
        const selectedOrganization =
          data?.dashboardView_Organizations?.content?.find(
            (e) => e.value === input.organizationId,
          );

        const subsidiaryOf = {
          organization: {
            id: input.organizationId,
            name: `${selectedOrganization?.label}`,
          },
        };
        queryClient.cancelQueries({ queryKey });

        queryClient.setQueryData<OrganizationQuery>(
          queryKey,
          (currentCache) => {
            return produce(currentCache, (draft) => {
              if (draft?.['organization']?.['subsidiaryOf']) {
                draft['organization']['subsidiaryOf'] = [subsidiaryOf];
              }
            });
          },
        );
        const previousEntries =
          queryClient.getQueryData<OrganizationQuery>(queryKey);
        const previousOrganizationsEntries = queryClient.getQueryData<
          InfiniteData<GetOrganizationsQuery>
        >(organizationsQueryKey);

        queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
          organizationsQueryKey,
          (currentCache) => {
            return produce(currentCache, (draft) => {
              const pageIndex =
                organizationsMeta.getOrganization.pagination.page - 1;
              const foundIndex = draft?.pages?.[
                pageIndex
              ]?.dashboardView_Organizations?.content?.findIndex(
                (o) => o.id === id,
              );

              if (typeof foundIndex === 'undefined' || foundIndex < 0) return;
              const dashboardContent =
                draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content;
              const item = dashboardContent?.[foundIndex];

              if (item && 'subsidiaryOf' in item) {
                item.subsidiaryOf = [subsidiaryOf];
              }
            });
          },
        );

        return { previousEntries, previousOrganizationsEntries };
      },
      onError: (_, __, context) => {
        queryClient.setQueryData(queryKey, context?.previousEntries);
        queryClient.setQueryData(
          organizationsQueryKey,
          context?.previousOrganizationsEntries,
        );
      },
      onSettled: () => {
        invalidateQuery();
        queryClient.invalidateQueries(organizationsQueryKey);
      },
    });
  const removeSubsidiaryToOrganizationMutation =
    useRemoveSubsidiaryToOrganizationMutation(client, {
      onMutate: (input) => {
        queryClient.cancelQueries({ queryKey });

        queryClient.setQueryData<OrganizationQuery>(
          queryKey,
          (currentCache) => {
            return produce(currentCache, (draft) => {
              if (draft?.['organization']?.['subsidiaryOf']) {
                draft['organization']['subsidiaryOf'] = [];
              }
            });
          },
        );
        const previousEntries =
          queryClient.getQueryData<OrganizationQuery>(queryKey);
        const previousOrganizationsEntries = queryClient.getQueryData<
          InfiniteData<GetOrganizationsQuery>
        >(organizationsQueryKey);
        queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
          organizationsQueryKey,
          (currentCache) => {
            return produce(currentCache, (draft) => {
              const pageIndex =
                organizationsMeta.getOrganization.pagination.page - 1;
              const foundIndex = draft?.pages?.[
                pageIndex
              ]?.dashboardView_Organizations?.content?.findIndex(
                (o) => o.id === id,
              );

              if (typeof foundIndex === 'undefined' || foundIndex < 0) return;

              const dashboardContent =
                draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content;
              if (
                dashboardContent &&
                dashboardContent[foundIndex] !== undefined
              ) {
                dashboardContent[foundIndex] = {
                  ...dashboardContent[foundIndex],
                  subsidiaryOf: [],
                };
              }
            });
          },
        );

        return { previousEntries, previousOrganizationsEntries };
      },

      onError: (_, __, context) => {
        queryClient.setQueryData(queryKey, context?.previousEntries);
        queryClient.setQueryData(
          organizationsQueryKey,
          context?.previousOrganizationsEntries,
        );
      },
      onSettled: () => {
        invalidateQuery();
        queryClient.invalidateQueries(organizationsQueryKey);
      },
    });

  const options = React.useMemo(() => {
    return (
      data?.dashboardView_Organizations?.content
        ?.filter((e) => !e.subsidiaryOf?.length && e.value !== id)
        .map((e) => ({
          label: e.label,
          value: e?.value,
        })) || []
    );
  }, [data?.dashboardView_Organizations?.content]);

  return (
    <Select
      isClearable
      value={parentOrg || ''}
      onChange={(e) => {
        if (!e && parentOrg) {
          removeSubsidiaryToOrganizationMutation.mutate({
            organizationId: parentOrg.value,
            subsidiaryId: id,
          });
        }
        if (e?.value) {
          addSubsidiaryToOrganizationMutation.mutate({
            input: {
              organizationId: e.value,
              subOrganizationId: id,
            },
          });
        }
      }}
      onInputChange={(inputValue) => setSearchTerm(inputValue)}
      isLoading={isLoading}
      options={options || []}
      placeholder='Parent organization'
      leftElement={<ArrowCircleBrokenUpLeft color='gray.500' mr='3' />}
    />
  );
};
