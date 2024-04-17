'use client';

import React, { FC } from 'react';
import { useParams } from 'next/navigation';

import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Flex } from '@ui/layout/Flex';
import { Contract, Organization } from '@graphql/types';
import { ContractCard } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractCard';
import { ARRForecast } from '@organization/src/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast';
import { ContractCard as NewContractCard } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractCard';
import { ContractModalsContextProvider } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractModalsContext';

import { Notes } from '../Notes';

interface ContractsProps {
  isLoading: boolean;
  organization?: Organization | null;
}
export const Contracts: FC<ContractsProps> = ({ isLoading, organization }) => {
  const id = useParams()?.id as string;
  const isNewContractUiEnabled = useFeatureIsOn('contract-new');

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
            <Flex
              key={`contract-card-${contract.metadata.id}`}
              flexDir='column'
              gap={4}
              w='full'
              mb={4}
            >
              {isNewContractUiEnabled ? (
                <ContractModalsContextProvider id={id}>
                  <NewContractCard
                    organizationId={id}
                    organizationName={organization?.name ?? ''}
                    data={(contract as Contract) ?? undefined}
                  />
                </ContractModalsContextProvider>
              ) : (
                <ContractCard
                  organizationId={id}
                  organizationName={organization?.name ?? ''}
                  data={(contract as Contract) ?? undefined}
                />
              )}
            </Flex>
          ))}
        </>
      )}

      <Notes id={id} data={organization} />
    </>
  );
};
