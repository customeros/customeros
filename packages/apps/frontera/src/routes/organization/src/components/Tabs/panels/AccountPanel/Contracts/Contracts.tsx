import { FC } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ContractCard } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractCard';
import { ARRForecast } from '@organization/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast.tsx';
import { ContractModalsContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext.tsx';
import { ContractModalStatusContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext.tsx';

import { Notes } from '../Notes';

interface ContractsProps {
  isLoading: boolean;
}
export const Contracts: FC<ContractsProps> = observer(({ isLoading }) => {
  const id = useParams()?.id as string;
  const store = useStore();
  const organizationStore = store.organizations.value.get(id);
  const contracts = organizationStore?.value.contracts?.map(
    (c) => c.metadata.id,
  );

  if (!organizationStore) return null;

  return (
    <>
      <ARRForecast
        renewalSunnary={
          organizationStore?.value?.accountDetails?.renewalSummary
        }
        name={organizationStore?.value?.name || ''}
        isInitialLoading={isLoading}
        currency={organizationStore?.value?.contracts?.[0]?.currency || 'USD'}
      />
      {contracts?.map((id) => {
        return (
          <div
            className='flex gap-4 flex-col w-full mb-4'
            key={`contract-card-${id}`}
          >
            <ContractModalStatusContextProvider id={id}>
              <ContractModalsContextProvider id={id}>
                <ContractCard
                  id={id}
                  organizationName={organizationStore?.value?.name || ''}
                />
              </ContractModalsContextProvider>
            </ContractModalStatusContextProvider>
          </div>
        );
      })}

      <Notes id={id} />
    </>
  );
});
