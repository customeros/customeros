'use client';

import React, { FC } from 'react';
import { useParams } from 'next/navigation';

import { Contract, Organization } from '@graphql/types';
import { ARRForecast } from '@organization/src/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast';
import { ContractCard as NewContractCard } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractCard';
import { ContractModalsContextProvider } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { ContractModalStatusContextProvider } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';
import { EditContractModalStoreContextProvider } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/EditContractModalStores';

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
                      <NewContractCard
                        organizationId={id}
                        organizationName={organization?.name ?? ''}
                        data={(contract as Contract) ?? undefined}
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
