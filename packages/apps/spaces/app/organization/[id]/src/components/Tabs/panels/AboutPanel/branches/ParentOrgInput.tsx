import React, { useState } from 'react';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';
import { Select } from '@ui/form/SyncSelect';
import { useAddSubsidiaryToOrganizationMutation } from '@organization/src/graphql/addSubsidiaryToOrganization.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useRemoveSubsidiaryToOrganizationMutation } from '@organization/src/graphql/removeSubsidiaryToOrganization.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/src/graphql/organization.generated';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';
import { produce } from 'immer';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { ComparisonOperator } from '@graphql/types';
import { useGetOrganizationOptionsQuery } from '@organization/src/graphql/getOrganizationOptions.generated';

interface ParentOrgInputProps {
  parentOrg: { label: string; value: string } | null;
  id: string;
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

  const options =
    data?.dashboardView_Organizations?.content
      ?.filter((e) => !e.subsidiaryOf?.length)
      .map((e) => ({
        label: e.label,
        value: e?.value,
      })) || null;

  return (
    <Select
      isClearable
      value={parentOrg || null}
      onChange={(e) => {
        if (!e && parentOrg) {
          removeSubsidiaryToOrganizationMutation.mutate({
            organizationId: parentOrg.value,
            subsidiaryId: id,
          });
        }
        addSubsidiaryToOrganizationMutation.mutate({
          input: {
            organizationId: e.value,
            subOrganizationId: id,
          },
        });
      }}
      onInputChange={(inputValue) => setSearchTerm(inputValue)}
      isLoading={isLoading}
      options={options || []}
      placeholder='Parent organization'
      leftElement={<ArrowCircleBrokenUpLeft color='gray.500' mr='3' />}
    />
  );
};
