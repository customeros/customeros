import React, { useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { toJS } from 'mobx';
import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Plus } from '@ui/media/icons/Plus';
import { Organization } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useCreateOrganizationMutation } from '@organizations/graphql/createOrganization.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/graphql/organization.generated';
import { useAddSubsidiaryToOrganizationMutation } from '@organization/graphql/addSubsidiaryToOrganization.generated';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

interface BranchesProps {
  id: string;
  isReadOnly?: boolean;
  branches?: Organization['subsidiaries'];
}

export const Branches: React.FC<BranchesProps> = ({
  id,
  isReadOnly,
  branches = [],
}) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const queryKey = useOrganizationQuery.getKey({ id });
  const organizationsQueryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );
  const invalidateQuery = () => queryClient.invalidateQueries({ queryKey });

  const addSubsidiaryToOrganizationMutation =
    useAddSubsidiaryToOrganizationMutation(client, {
      onMutate: () => {
        const subsidiaryOf = {
          organization: {
            id,
            name: '',
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
                (o) => o.metadata.id === id,
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

      onSuccess: (_data, variables, _context) => {
        navigate(`/organization/${variables?.input?.subsidiaryId}`);

        toastSuccess(
          'Organization created',
          `${variables?.input?.subsidiaryId}-created`,
        );
      },
      onError: (_, __, context) => {
        queryClient.setQueryData(queryKey, context?.previousEntries);
        queryClient.setQueryData(
          organizationsQueryKey,
          context?.previousOrganizationsEntries,
        );
        toastError(
          'Failed to create organization',
          'create-new-organization-error',
        );
      },
      onSettled: () => {
        invalidateQuery();
        queryClient.invalidateQueries({ queryKey: organizationsQueryKey });
      },
    });
  const createOrganization = useCreateOrganizationMutation(client, {
    onMutate: () => {
      const pageIndex = organizationsMeta.getOrganization.pagination.page - 1;
      queryClient.cancelQueries({ queryKey });

      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            if (!draft) return;

            const page = draft.pages?.[pageIndex];
            let content = page?.dashboardView_Organizations?.content;

            const emptyRow = produce(content?.[0], (draft) => {
              if (!draft) return;
              draft.metadata.id = Math.random().toString();
              draft.name = '';
              draft.website = '';
              draft.owner = null;
              draft.accountDetails = null;
            });

            if (!emptyRow) return;
            content = [emptyRow, ...(content ?? [])];
          });
        },
      );

      setOrganizationsMeta((prev) =>
        produce(prev, (draft) => {
          draft.getOrganization.pagination.page = 1;
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        context?.previousEntries,
      );
      toastError('Failed to create organization', 'create-organization-error');
    },
    onSuccess: ({ organization_Create: { id: createdOrgId } }) => {
      addSubsidiaryToOrganizationMutation.mutate({
        input: {
          organizationId: id,
          subsidiaryId: createdOrgId,
        },
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });
  const handleCreateOrganization = useCallback(() => {
    createOrganization.mutate({ input: { name: 'Unnamed' } });
  }, [createOrganization]);

  console.log(toJS(branches));

  return (
    <Card className='w-full mt-2 p-4 bg-white rounded-md border-1 shadow-lg'>
      <CardHeader className='flex mb-4 items-center justify-between'>
        <h2 className='text-base'>Branches</h2>
        {!isReadOnly && (
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Add'
            onClick={handleCreateOrganization}
            icon={<Plus className='size-4' />}
          />
        )}
      </CardHeader>
      <CardContent className='flex flex-col p-0 pt-0 gap-2 items-baseline'>
        {/* {branches?.map(({ organization }) =>
          organization?.metadata?.id ? (
            <Link
              className='line-clamp-1 break-keep text-gray-700 hover:text-primary-600 no-underline hover:underline'
              to={`/organization/${organization.metadata?.id}?tab=about`}
              key={`subsidiaries-${organization.metadata?.id}`}
            >
              {organization?.name || 'Unknown'}
            </Link>
          ) : null,
        )} */}
      </CardContent>
    </Card>
  );
};
