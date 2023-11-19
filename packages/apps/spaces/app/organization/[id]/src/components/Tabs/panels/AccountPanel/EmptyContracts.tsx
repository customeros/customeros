'use client';

import { FC } from 'react';
import { useParams } from 'next/navigation';

import { produce } from 'immer';
import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Star06 } from '@ui/media/icons/Star06';
import { Heading } from '@ui/typography/Heading';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useCreateContractMutation } from '@organization/src/graphql/createContract.generated';
import {
  User,
  DataSource,
  ContractStatus,
  ContractRenewalCycle,
} from '@graphql/types';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const EmptyContracts: FC<{ name: string }> = ({ name }) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data: session } = useSession();
  const id = useParams()?.id as string;
  const queryKey = useGetContractsQuery.getKey({ id });

  const createContract = useCreateContractMutation(client, {
    onMutate: () => {
      const contract = {
        appSource: DataSource.Openline,
        contractUrl: '',
        createdAt: new Date().toISOString(),
        createdBy: [session?.user] as unknown as User,
        externalLinks: [],
        renewalCycle: ContractRenewalCycle.None,
        id: `created-contract-${Math.random().toString()}`,
        name: '',
        owner: null,
        source: DataSource.Openline,
        sourceOfTruth: DataSource.Openline,
        status: ContractStatus.Draft,
        updatedAt: new Date().toISOString(),
        serviceLineItems: [],
      };
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts'] = [contract];
          }
        });
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },

    onSuccess: (_, variables) => {
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
      queryClient.invalidateQueries(queryKey);
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
