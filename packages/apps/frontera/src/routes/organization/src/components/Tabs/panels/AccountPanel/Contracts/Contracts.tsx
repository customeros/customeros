import { FC } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ContractCard } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractCard';
import { ARRForecast } from '@organization/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast.tsx';
import { ContractModalsContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext.tsx';
import { ContractModalStatusContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext.tsx';
import { EditContractModalStoreContextProvider } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/EditContractModalStores.tsx';

import { Notes } from '../Notes';

interface ContractsProps {
  isLoading: boolean;
}
export const Contracts: FC<ContractsProps> = observer(({ isLoading }) => {
  const id = useParams()?.id as string;
  const store = useStore();
  const organizationStore = store.organizations.value.get(id)?.value;

  const contracts = organizationStore?.contracts;

  return (
    <>
      <>
        <ARRForecast
          renewalSunnary={organizationStore?.accountDetails?.renewalSummary}
          name={organizationStore?.name || ''}
          isInitialLoading={isLoading}
          currency={organizationStore?.contracts?.[0]?.currency || 'USD'}
        />
        {contracts?.map((c) => {
          return (
            <div
              className='flex gap-4 flex-col w-full mb-4'
              key={`contract-card-${c.metadata.id}`}
            >
              <ContractModalStatusContextProvider id={c.metadata.id}>
                <EditContractModalStoreContextProvider>
                  <ContractModalsContextProvider id={c.metadata.id}>
                    <ContractCard
                      values={c}
                      organizationName={organizationStore?.name || ''}
                    />
                  </ContractModalsContextProvider>
                </EditContractModalStoreContextProvider>
              </ContractModalStatusContextProvider>
            </div>
          );
        })}
      </>

      <Notes id={id} data={organizationStore} />
    </>
  );
});
