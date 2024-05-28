import { FC } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Organization } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { ContractCard } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractCard';
import { ARRForecast } from '@organization/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast.tsx';
import { ContractModalsContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext.tsx';
import { ContractModalStatusContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext.tsx';
import { EditContractModalStoreContextProvider } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/EditContractModalStores.tsx';

import { Notes } from '../Notes';

interface ContractsProps {
  isLoading: boolean;
  organization?: Organization | null;
}
export const Contracts: FC<ContractsProps> = observer(
  ({ isLoading, organization }) => {
    const id = useParams()?.id as string;
    const store = useStore();
    const contracts = store.organizations.value
      .get(id)
      ?.value?.contracts?.map((e) => e.metadata.id);

    return (
      <>
        {!!organization?.contracts && (
          <>
            <ARRForecast
              renewalSunnary={organization?.accountDetails?.renewalSummary}
              name={organization?.name || ''}
              isInitialLoading={isLoading}
              currency={organization?.contracts?.[0]?.currency || 'USD'}
            />
            {contracts?.map((id) => {
              return (
                <div
                  className='flex gap-4 flex-col w-full mb-4'
                  key={`contract-card-${id}`}
                >
                  <ContractModalStatusContextProvider
                    id={id}
                    upcomingInvoices={[]}
                    nextInvoice={undefined}
                    committedPeriodInMonths={0}
                  >
                    <EditContractModalStoreContextProvider>
                      <ContractModalsContextProvider id={id}>
                        <ContractCard
                          organizationId={id}
                          contractId={id}
                          organizationName={organization?.name || ''}
                        />
                      </ContractModalsContextProvider>
                    </EditContractModalStoreContextProvider>
                  </ContractModalStatusContextProvider>
                </div>
              );
            })}
          </>
        )}

        <Notes id={id} data={organization} />
      </>
    );
  },
);
