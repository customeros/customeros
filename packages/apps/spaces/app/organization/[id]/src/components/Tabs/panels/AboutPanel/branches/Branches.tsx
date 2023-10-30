import React, { useCallback } from 'react';
import { Card } from '@ui/layout/Card';
import { CardHeader, Heading, IconButton, VStack } from '@chakra-ui/react';
import { Plus } from '@ui/media/icons/Plus';
import { CardBody } from '@chakra-ui/card';
import { Link } from '@ui/navigation/Link';

import { Organization } from '@graphql/types';
import { useAddSubsidiaryToOrganizationMutation } from '@organization/src/graphql/addSubsidiaryToOrganization.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/src/graphql/organization.generated';
import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useCreateOrganizationMutation } from '@organizations/graphql/createOrganization.generated';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useRouter } from 'next/navigation';

interface BranchesProps {
  branches?: Organization['subsidiaries'];
  id: string;
}

export const Branches: React.FC<BranchesProps> = ({ id, branches = [] }) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { push } = useRouter();
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();
  const queryKey = useOrganizationQuery.getKey({ id });
  const organizationsQueryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );
  const invalidateQuery = () => queryClient.invalidateQueries(queryKey);

  const addSubsidiaryToOrganizationMutation =
    useAddSubsidiaryToOrganizationMutation(client, {
      onMutate: ({ input }) => {
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

      onSuccess: (data, variables, context) => {
        push(`/organization/${variables?.input?.subOrganizationId}`);

        toastSuccess(
          'Organization created',
          `${variables?.input?.subOrganizationId}-created`,
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
        queryClient.invalidateQueries(organizationsQueryKey);
      },
    });
  const createOrganization = useCreateOrganizationMutation(client, {
    onMutate: () => {
      const pageIndex = organizationsMeta.getOrganization.pagination.page - 1;
      queryClient.cancelQueries(queryKey);

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
              draft.id = Math.random().toString();
              draft.name = '';
              draft.website = '';
              draft.owner = null;
              draft.lastTouchPointTimelineEvent = null;
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
          subOrganizationId: createdOrgId,
        },
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries(queryKey);
    },
  });
  const handleCreateOrganization = useCallback(() => {
    createOrganization.mutate({ input: { name: '' } });
  }, [createOrganization]);

  return (
    <Card size='sm' width='full' mt={2}>
      <CardHeader
        display='flex'
        alignItems='center'
        justifyContent='space-between'
        pb={4}
      >
        <Heading fontSize={'md'}>Branches</Heading>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Add'
          onClick={handleCreateOrganization}
          icon={<Plus boxSize='4' />}
        />
      </CardHeader>
      <CardBody as={VStack} pt={0} gap={2} alignItems='baseline'>
        {branches.map(({ organization }) =>
          organization?.id ? (
            <Link
              noOfLines={1}
              wordBreak='keep-all'
              href={`/organization/${organization.id}?tab=about`}
              key={`subsidiaries-${organization.id}`}
              color='gray.700'
              _hover={{ color: 'primary.600' }}
            >
              {organization?.name || 'Unknown'}
            </Link>
          ) : null,
        )}
      </CardBody>
    </Card>
  );
};
