'use client';

import React from 'react';

import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';

import { BankAccount } from '@graphql/types';
import { FormSwitch } from '@ui/form/Switch/FormSwitch2';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { AddAccountButton } from './AddAccountButton';
import { BankTransferCard } from './BankTransferCard';

export const BankTransferAccountList = ({ formId }: { formId: string }) => {
  const client = getGraphQLClient();
  const { data } = useBankAccountsQuery(client);

  const existingAccountCurrencies = data?.bankAccounts.map(
    (account) => account.currency as string,
  );

  return (
    <>
      <div className='flex justify-between items-center'>
        <span className='text-sm whitespace-nowrap font-semibold'>
          Bank transfer
        </span>
        <div className='flex items-center'>
          <AddAccountButton
            existingCurrencies={existingAccountCurrencies ?? []}
          />
          <div>
            <FormSwitch
              name='canPayWithBankTransfer'
              formId={formId}
              size='sm'
            />
          </div>
        </div>
      </div>

      {data?.bankAccounts?.map((account) => (
        <React.Fragment key={account.metadata.id}>
          <BankTransferCard
            account={account as BankAccount}
            existingCurrencies={existingAccountCurrencies ?? []}
          />
        </React.Fragment>
      ))}
    </>
  );
};
