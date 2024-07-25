import { FC, PropsWithChildren } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

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
import { OrganizationPanel } from '@organization/components/Tabs';
import { Contracts } from '@organization/components/Tabs/panels/AccountPanel/Contracts/Contracts';
import { RelationshipButton } from '@organization/components/Tabs/panels/AccountPanel/RelationshipButton';

import { Notes } from './Notes';
import { EmptyContracts } from './EmptyContracts';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import {
  useAccountPanelStateContext,
  AccountModalsContextProvider,
} from './context/AccountModalsContext';

const AccountPanelComponent = observer(() => {
  const navigate = useNavigate();
  const store = useStore();
  const baseCurrency = store.settings.tenant.value?.baseCurrency;

  const id = useParams()?.id as string;

  const { isModalOpen } = useAccountPanelStateContext();

  const organizationStore = store.organizations.value.get(id);
  const organization = organizationStore?.value;

  if (store.organizations.isLoading) {
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
      currency: baseCurrency || Currency.Usd,
      name: `${
        organization?.name?.length ? `${organization?.name}'s` : "Unnamed's"
      } contract`,
    });
  };

  const isCreating = organization?.metadata?.id
    ? Boolean(store.contracts.isPending.get(organization?.metadata?.id))
    : false;

  if (!organizationStore?.contracts?.length) {
    return (
      <EmptyContracts isPending={isCreating} onCreate={handleCreate}>
        <Notes id={id} />
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
              {store.invoices.isLoading ? (
                <Skeleton className='h-3 w-3 ml-1' />
              ) : (
                organizationStore.invoices?.length
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
                data-Test='org-account-nonempty-new-contract'
                onClick={() => handleCreate()}
              />
            </Tooltip>

            <RelationshipButton />
          </div>
        }
        shouldBlockPanelScroll={isModalOpen}
      >
        <Contracts isLoading={isCreating} />
      </OrganizationPanel>
    </>
  );
});

export const AccountPanel: FC<PropsWithChildren> = () => (
  <AccountModalsContextProvider>
    <AccountPanelComponent />
  </AccountModalsContextProvider>
);
