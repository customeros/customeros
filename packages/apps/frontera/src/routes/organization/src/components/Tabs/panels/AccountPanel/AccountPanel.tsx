import { FC, PropsWithChildren } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { toJS } from 'mobx';
import { observer } from 'mobx-react-lite';
import { useBaseCurrencyQuery } from '@settings/graphql/getBaseCurrency.generated';

import { Currency } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus';
import { DateTimeUtils } from '@utils/date.ts';
import { Button } from '@ui/form/Button/Button';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractsQuery } from '@organization/graphql/getContracts.generated';
import { useGetInvoicesCountQuery } from '@organization/graphql/getInvoicesCount.generated';
import { Contracts } from '@organization/components/Tabs/panels/AccountPanel/Contracts/Contracts';
import { RelationshipButton } from '@organization/components/Tabs/panels/AccountPanel/RelationshipButton';

import { Notes } from './Notes';
import { EmptyContracts } from './EmptyContracts';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import {
  useAccountPanelStateContext,
  AccountModalsContextProvider,
} from './context/AccountModalsContext';

const AccountPanelComponent = observer(() => {
  const navigate = useNavigate();
  const store = useStore();
  const client = getGraphQLClient();

  const id = useParams()?.id as string;

  const { isModalOpen } = useAccountPanelStateContext();
  const { data, isLoading } = useGetContractsQuery(client, {
    id,
  });
  const { data: invoicesCountData, isFetching: isFetchingInvoicesCount } =
    useGetInvoicesCountQuery(client, {
      organizationId: id,
    });
  const { data: baseCurrencyData } = useBaseCurrencyQuery(client);
  const organizationStore = store.organizations.value.get(id)?.value;
  if (isLoading) {
    return <AccountPanelSkeleton />;
  }

  const handleCreate = () => {
    store.contracts.create({
      organizationId: id,
      serviceStarted: DateTimeUtils.addDays(
        new Date().toString(),
        1,
      ).toISOString(),
      committedPeriodInMonths: 1,
      currency: baseCurrencyData?.tenantSettings?.baseCurrency || Currency.Usd,
      name: `${
        data?.organization?.name?.length
          ? `${data?.organization?.name}'s`
          : "Unnamed's"
      } contract`,
    });
  };

  if (!organizationStore?.contracts?.length) {
    return (
      <EmptyContracts
        isPending={store.contracts.isLoading}
        onCreate={handleCreate}
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
                isLoading={store.contracts.isLoading}
                isDisabled={store.contracts.isLoading}
                icon={
                  store.contracts.isLoading ? (
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
        <Contracts isLoading={isLoading} />
      </OrganizationPanel>
    </>
  );
});

export const AccountPanel: FC<PropsWithChildren> = () => (
  <AccountModalsContextProvider>
    <AccountPanelComponent />
  </AccountModalsContextProvider>
);
