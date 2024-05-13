import { FC } from 'react';
import { useParams } from 'react-router-dom';

import { Contract, Organization } from '@graphql/types';
import { ARRForecast } from '@organization/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast';
import { ContractCard as NewContractCard } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractCard';
import { ContractModalsContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { ContractModalStatusContextProvider } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

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
          {organization?.contracts.map((contract) => (
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
                <ContractModalsContextProvider id={id}>
                  <NewContractCard
                    organizationId={id}
                    organizationName={organization?.name ?? ''}
                    data={(contract as Contract) ?? undefined}
                  />
                </ContractModalsContextProvider>
              </ContractModalStatusContextProvider>
            </div>
          ))}
        </>
      )}

      <Notes id={id} data={organization} />
    </>
  );
};
