'use client';

import React, { FC } from 'react';
import { useParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Contract, Organization } from '@graphql/types';
import { ARRForecast } from '@organization/src/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast';

import { Notes } from '../Notes';
import { ContractCard } from '../Contract/ContractCard';

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
            currency={organization?.contracts[0]?.currency}
          />
          {organization?.contracts.map((contract) => (
            <Flex
              key={`contract-card-${contract.id}`}
              flexDir='column'
              gap={4}
              mb={4}
            >
              <ContractCard
                organizationId={id}
                organizationName={organization?.name ?? ''}
                data={(contract as Contract) ?? undefined}
              />
            </Flex>
          ))}
        </>
      )}

      <Notes id={id} data={organization} />
    </>
  );
};
