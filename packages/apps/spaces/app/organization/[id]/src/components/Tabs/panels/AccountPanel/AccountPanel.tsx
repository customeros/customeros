'use client';

import React, { FC, PropsWithChildren } from 'react';
import { useParams, useRouter } from 'next/navigation';

import { produce } from 'immer';
import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';
import { useBaseCurrencyQuery } from '@settings/graphql/getBaseCurrency.generated';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Spinner } from '@ui/feedback/Spinner';
import { IconButton } from '@ui/form/IconButton';
import { toastError } from '@ui/presentation/Toast';
import { Skeleton } from '@ui/presentation/Skeleton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCreateContractMutation } from '@organization/src/graphql/createContract.generated';
import { useGetInvoicesCountQuery } from '@organization/src/graphql/getInvoicesCount.generated';
import { Contracts } from '@organization/src/components/Tabs/panels/AccountPanel/Contracts/Contracts';
import { RelationshipButton } from '@organization/src/components/Tabs/panels/AccountPanel/RelationshipButton';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import {
  User,
  Currency,
  DataSource,
  Organization,
  ContractStatus,
  ContractRenewalCycle,
} from '@graphql/types';

import { Notes } from './Notes';
import { EmptyContracts } from './EmptyContracts';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import {
  useAccountPanelStateContext,
  AccountModalsContextProvider,
} from './context/AccountModalsContext';

const AccountPanelComponent = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data: session } = useSession();

  const id = useParams()?.id as string;
  const router = useRouter();
  const queryKey = useGetContractsQuery.getKey({ id });

  const { isModalOpen } = useAccountPanelStateContext();
  const { data, isLoading } = useGetContractsQuery(client, {
    id,
  });

  const { data: invoicesCountData, isFetching: isFetchingInvoicesCount } =
    useGetInvoicesCountQuery(client, {
      organizationId: id,
    });
  const { data: baseCurrencyData } = useBaseCurrencyQuery(client);

  const createContract = useCreateContractMutation(client, {
    onMutate: () => {
      const contract = {
        contractUrl: '',
        metadata: {
          id: `created-contract-${Math.random().toString()}`,
          created: new Date().toISOString(),
          lastUpdated: new Date().toISOString(),

          source: DataSource.Openline,
        },
        createdBy: [session?.user] as unknown as User,
        externalLinks: [],
        contractRenewalCycle: ContractRenewalCycle.MonthlyRenewal,
        contractName: `${
          data?.organization?.name?.length
            ? `${data?.organization?.name}'s`
            : "Unnamed's"
        } contract`,
        owner: null,
        autoRenew: false,
        contractStatus: ContractStatus.Draft,
        contractLineItems: [],
        billingEnabled: false,
        approved: false,
        upcomingInvoices: [],
      };
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts'] = [
              ...(currentCache?.organization?.contracts || []),
              contract,
            ];
          }
        });
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },

    onError: (_, __, context) => {
      queryClient.setQueryData(queryKey, context?.previousEntries);
      toastError(
        'Failed to create contract',
        'create-new-contract-for-organization-error',
      );
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });

  if (isLoading) {
    return <AccountPanelSkeleton />;
  }

  if (!data?.organization?.contracts?.length) {
    return (
      <EmptyContracts
        name={data?.organization?.name || ''}
        baseCurrency={
          baseCurrencyData?.tenantSettings?.baseCurrency || Currency.Usd
        }
      >
        <Notes id={id} data={data?.organization} />
      </EmptyContracts>
    );
  }

  return (
    <>
      <OrganizationPanel
        title='Account'
        withFade
        bottomActionItem={
          <Button
            borderRadius={0}
            bg='gray.25'
            p={7}
            justifyContent='space-between'
            alignItems='center'
            rightIcon={<ChevronRight boxSize={4} color='gray.400' />}
            variant='ghost'
            _hover={{
              bg: 'gray.25',
              '& svg': {
                color: 'gray.500',
              },
            }}
            onClick={() => router.push(`?tab=invoices`)}
          >
            <Text
              fontSize='sm'
              fontWeight='semibold'
              display='inline-flex'
              alignItems='center'
            >
              Invoices •{' '}
              {isFetchingInvoicesCount ? (
                <Skeleton as='span' height={3} width={2} ml={1} />
              ) : (
                invoicesCountData?.invoices.totalElements
              )}
            </Text>
          </Button>
        }
        actionItem={
          <Flex alignItems='center'>
            <Tooltip label='Create new contract'>
              <IconButton
                color='gray.500'
                mr={1}
                variant='ghost'
                isLoading={createContract.isPending}
                isDisabled={createContract.isPending}
                icon={createContract.isPending ? <Spinner /> : <Plus />}
                size='xs'
                aria-label='Create new contract'
                onClick={() =>
                  createContract.mutate({
                    input: {
                      organizationId: id,
                      contractRenewalCycle: ContractRenewalCycle.MonthlyRenewal,
                      currency:
                        baseCurrencyData?.tenantSettings?.baseCurrency ||
                        Currency.Usd,
                      name: `${
                        data?.organization?.name?.length
                          ? `${data?.organization?.name}'s`
                          : "Unnamed's"
                      } contract`,
                    },
                  })
                }
              />
            </Tooltip>

            <RelationshipButton />
          </Flex>
        }
        shouldBlockPanelScroll={isModalOpen}
      >
        <Contracts
          isLoading={isLoading}
          organization={data?.organization as Organization}
        />
      </OrganizationPanel>
    </>
  );
};

export const AccountPanel: FC<PropsWithChildren> = () => (
  <AccountModalsContextProvider>
    <AccountPanelComponent />
  </AccountModalsContextProvider>
);
