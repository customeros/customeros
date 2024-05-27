import { FC } from 'react';
import { useParams } from 'react-router-dom';

import { Contract, Organization } from '@graphql/types';
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
export const Contracts: FC<ContractsProps> = ({ isLoading, organization }) => {
  const id = useParams()?.id as string;

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
          {organization?.contracts.map((contract) => {
            return (
              <div
                className='flex gap-4 flex-col w-full mb-4'
                key={`contract-card-${contract.metadata.id}`}
              >
                <ContractModalStatusContextProvider
                  id={id}
                  upcomingInvoices={contract?.upcomingInvoices}
                  nextInvoice={contract?.billingDetails?.nextInvoicing}
                  committedPeriodInMonths={contract?.committedPeriodInMonths}
                >
                  <EditContractModalStoreContextProvider>
                    <ContractModalsContextProvider id={contract.metadata.id}>
                      <ContractCard
                        organizationId={id}
                        organizationName={organization?.name || ''}
                        data={contract as Contract}
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
};
