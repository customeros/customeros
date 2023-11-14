'use client';

import { FC, useRef } from 'react';
import { useParams } from 'next/navigation';

import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';

import { User } from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Star06 } from '@ui/media/icons/Star06';
import { Heading } from '@ui/typography/Heading';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { NEW_DATE } from '@organization/src/components/Timeline/OrganizationTimeline';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { useCreateContractMutation } from '@organization/src/graphql/createContract.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/src/graphql/organization.generated';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const EmptyContracts: FC<{ name: string }> = ({ name }) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data: session } = useSession();
  const id = useParams()?.id as string;
  const queryKey = useOrganizationQuery.getKey({ id });
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const createContract = useCreateContractMutation(client, {
    onMutate: ({ input }) => {
      // todo uncomment when contract field is added to organization query
      // const contract = {
      //   __typename: 'Contract',
      //   appSource: '',
      //   contractUrl: '',
      //   createdAt: new Date().toISOString(),
      //   createdBy: [session?.user] as unknown as User,
      //   endedAt: new Date().toISOString(),
      //   externalLinks: [],
      //   id: 'abcd1234',
      //   name: '',
      //   owner: null,
      //   renewalCycle: ContractRenewalCycle.None,
      //   serviceStartedAt: new Date().toISOString(),
      //   signedAt: new Date().toISOString(),
      //   source: DataSource.Openline,
      //   sourceOfTruth: DataSource.Openline,
      //   status: ContractStatus.Draft,
      //   updatedAt: new Date().toISOString(),
      // };
      queryClient.cancelQueries({ queryKey });
      // todo uncomment when organization query will contain contracts
      // queryClient.setQueryData<OrganizationQuery>(queryKey, (currentCache) => {
      //   return produce(currentCache, (draft) => {
      //     if (draft?.['organization']?.['contracts']) {
      //       draft['organization']['contracts'] = [contract];
      //     }
      //   });
      // });
      const previousEntries =
        queryClient.getQueryData<OrganizationQuery>(queryKey);

      return { previousEntries };
    },

    onSuccess: (data, variables, context) => {
      queryClient.setQueryData(
        useInfiniteGetTimelineQuery.getKey({
          organizationId: id,
          from: NEW_DATE,
          size: 50,
        }),
        (oldData) => {
          const newEvent = {
            __typename: 'Action',
            id: `timeline-event-action-new-id-${new Date()}`,
            actionType: 'CONTRACT_UPDATED',
            appSource: 'customer-os-api',
            createdAt: new Date().toISOString(),
            actionCreatedBy: [session?.user] as unknown as User,
            content: `Contract created`,
          };

          // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
          if (!oldData || !oldData.pages?.length) {
            return {
              pages: [
                {
                  organization: {
                    id,
                    timelineEventsTotalCount: 1,
                    timelineEvents: [newEvent],
                  },
                },
              ],
            };
          }

          // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
          const firstPage = oldData.pages[0] ?? {};
          // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
          const pages = oldData.pages?.slice(1);

          const firstPageWithEvent = {
            ...firstPage,
            organization: {
              ...firstPage?.organization,
              timelineEvents: [
                ...(firstPage?.organization?.timelineEvents ?? []),
                newEvent,
              ],
              timelineEventsTotalCount:
                (firstPage?.organization?.timelineEventsTotalCount ?? 0) + 1,
            },
          };

          return {
            ...oldData,
            pages: [firstPageWithEvent, ...pages],
          };
        },
      );

      toastSuccess(
        'Contract created',
        `${variables?.input?.organizationId}-contract-created`,
      );
    },
    onError: (_, __, context) => {
      queryClient.setQueryData(queryKey, context?.previousEntries);
      toastError(
        'Failed to create contract',
        'create-new-contract-for-organization-error',
      );
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
      }, 1000);
    },
  });

  return (
    <OrganizationPanel
      title='Account'
      actionItem={
        <Button
          size='xs'
          variant='outline'
          type='button'
          isDisabled
          borderRadius='16px'
        >
          Prospect
        </Button>
      }
    >
      <Flex
        mt={4}
        w='full'
        boxShadow={'none'}
        flexDir='column'
        justifyItems='center'
        alignItems='center'
      >
        <FeaturedIcon colorScheme='primary' mb={2} size='lg'>
          <Star06 boxSize={4} />
        </FeaturedIcon>
        <Heading mb={1} size='sm' fontWeight='semibold'>
          Create new contract
        </Heading>
        <Text fontSize='sm'>
          Create new contract for
          <Text as='span' fontWeight='medium' ml={1}>
            {name}
          </Text>
        </Text>
        <Button
          fontSize='sm'
          size='sm'
          colorScheme='primary'
          mt={6}
          variant='outline'
          width='fit-content'
          onClick={() =>
            createContract.mutate({
              input: { organizationId: id },
            })
          }
        >
          New contract
        </Button>
      </Flex>
    </OrganizationPanel>
  );
};
