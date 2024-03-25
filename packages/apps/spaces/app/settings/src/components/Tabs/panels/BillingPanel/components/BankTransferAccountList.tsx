'use client';

import React from 'react';

import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';

import { BankAccount } from '@graphql/types';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { AddAccountButton } from './AddAccountButton';
import { BankTransferCard } from './BankTransferCard';

export const BankTransferAccountList = ({
  formId,
  organizationName,
}: {
  formId: string;
  organizationName?: string | null;
}) => {
  const client = getGraphQLClient();
  const { data } = useBankAccountsQuery(client);

  const existingAccountCurrencies = data?.bankAccounts.map(
    (account) => account.currency as string,
  );

  return (
    <>
      <div className='flex items-center justify-between'>
        <span className='text-sm font-semibold whitespace-nowrap'>
          Bank transfer
        </span>
        <span className='flex items-center'>
          <AddAccountButton
            existingCurrencies={existingAccountCurrencies ?? []}
            organizationName={organizationName}
          />
          <div>
            <FormSwitch
              name='canPayWithBankTransfer'
              formId={formId}
              size='sm'
            />
          </div>
        </span>
      </div>

      {data?.bankAccounts?.map((account) => (
        <React.Fragment key={account.metadata.id}>
          <BankTransferCard account={account as BankAccount} />
        </React.Fragment>
      ))}
    </>
  );
};
