import { FC, PropsWithChildren } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { useQueryClient } from '@tanstack/react-query';
import { useBaseCurrencyQuery } from '@settings/graphql/getBaseCurrency.generated';

import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCreateContractMutation } from '@organization/graphql/createContract.generated';
import { useGetInvoicesCountQuery } from '@organization/graphql/getInvoicesCount.generated';
import { Contracts } from '@organization/components/Tabs/panels/AccountPanel/Contracts/Contracts';
import { RelationshipButton } from '@organization/components/Tabs/panels/AccountPanel/RelationshipButton';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
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
  const navigate = useNavigate();
  const store = useStore();
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const id = useParams()?.id as string;
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
        createdBy: [store.session?.value] as unknown as User,
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
      // @ts-expect-error - will be removed when we go to stores
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        if (!currentCache) {
          return {
            organization: {
              contracts: [contract],
            },
          };
        }

        return {
          ...currentCache,
          organization: {
            ...currentCache?.organization,
            contracts: [
              contract,
              ...(currentCache?.organization?.contracts || []),
            ],
          },
        };
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },

    onError: () => {
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
  const handleCreate = () => {
    createContract.mutate({
      input: {
        organizationId: id,
        contractRenewalCycle: ContractRenewalCycle.MonthlyRenewal,
        currency:
          baseCurrencyData?.tenantSettings?.baseCurrency || Currency.Usd,
        name: `${
          data?.organization?.name?.length
            ? `${data?.organization?.name}'s`
            : "Unnamed's"
        } contract`,
      },
    });
  };

  return (
    <>
      <OrganizationPanel
        title='Account'
        withFade
        bottomActionItem={
          <Button
            className='rounded-none bg-gray-25 p-7 justify-between items-center hover:bg-gray-25 group'
            rightIcon={
              <ChevronRight className='size-4 text-gray-400 group-hover:text-gray-500' />
            }
            variant='ghost'
            onClick={() => navigate(`?tab=invoices`)}
          >
            <p className='text-sm font-semibold inline-flex items-center'>
              Invoices •{' '}
              {isFetchingInvoicesCount ? (
                <Skeleton className='h-3 w-3 ml-1' />
              ) : (
                invoicesCountData?.invoices.totalElements
              )}
            </p>
          </Button>
        }
        actionItem={
          <div className='flex items-center'>
            <Tooltip label='Create new contract'>
              <IconButton
                className='text-gray-500 mr-1'
                variant='ghost'
                isLoading={createContract.isPending}
                isDisabled={createContract.isPending}
                icon={
                  createContract.isPending ? (
                    <Spinner
                      className='text-gray-500 fill-gray-700'
                      size='sm'
                      label='Creating contract...'
                    />
                  ) : (
                    <Plus />
                  )
                }
                size='xs'
                aria-label='Create new contract'
                onClick={() => handleCreate()}
              />
            </Tooltip>

            <RelationshipButton />
          </div>
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
