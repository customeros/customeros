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
  const contracts = organizationStore?.contracts;

  return (
    <>
      <ARRForecast
        isInitialLoading={isLoading}
        name={organizationStore?.value.name || ''}
        currency={contracts?.[0]?.currency || 'USD'}
        renewalSunnary={organizationStore?.value.accountDetails?.renewalSummary}
      />
      {contracts?.map((c) => {
        return (
          <div
            key={`contract-card-${c.metadata.id}`}
            className='flex gap-4 flex-col w-full mb-4'
          >
            <ContractModalStatusContextProvider id={c.metadata.id}>
              <ContractModalsContextProvider id={c.metadata.id}>
                <ContractCard
                  values={c}
                  organizationName={organizationStore?.value.name || ''}
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
