import { useParams } from 'react-router-dom';
import { FC, PropsWithChildren } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { File02 } from '@ui/media/icons/File02';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCreateContractMutation } from '@organization/graphql/createContract.generated';
import { RelationshipButton } from '@organization/components/Tabs/panels/AccountPanel/RelationshipButton';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import {
  User,
  Currency,
  DataSource,
  ContractStatus,
  ContractRenewalCycle,
} from '@graphql/types';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const EmptyContracts: FC<
  PropsWithChildren<{ name: string; baseCurrency: Currency }>
> = ({ name, baseCurrency, children }) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const store = useStore();
  const id = useParams()?.id as string;
  const queryKey = useGetContractsQuery.getKey({ id });

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
        createdBy: store.session.value as unknown as User,
        externalLinks: [],
        contractRenewalCycle: ContractRenewalCycle.None,
        contractName: `${name?.length ? `${name}'s` : "Unnamed's"} contract`,
        owner: null,
        contractStatus: ContractStatus.Draft,
        contractLineItems: [],
        billingEnabled: false,
        autoRenew: false,
        approved: false,
        upcomingInvoices: [],
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

  return (
    <OrganizationPanel title='Account' actionItem={<RelationshipButton />}>
      <article className='my-4 w-full flex flex-col items-center'>
        <FeaturedIcon className='mb-4' colorScheme='primary' size='lg'>
          <File02 className='size-4' />
        </FeaturedIcon>
        <h1 className='text-md font-semibold'>Draft a new contract</h1>

        <Button
          size='sm'
          className='text-sm mt-6 w-fit'
          isLoading={createContract.isPending}
          isDisabled={createContract.isPending}
          colorScheme='primary'
          variant='outline'
          onClick={() =>
            createContract.mutate({
              input: {
                organizationId: id,
                currency: baseCurrency,
                name: `${name?.length ? `${name}'s` : "Unnamed's"} contract`,
                contractRenewalCycle: ContractRenewalCycle.MonthlyRenewal,
                serviceStarted: DateTimeUtils.addDays(
                  new Date().toISOString(),
                  1,
                ),
              },
            })
          }
        >
          {createContract.isPending ? 'Creating contract...' : 'New contract'}
        </Button>
      </article>
      {children}
    </OrganizationPanel>
  );
};
