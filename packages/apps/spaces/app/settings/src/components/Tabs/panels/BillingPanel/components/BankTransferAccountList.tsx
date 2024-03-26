'use client';

import React from 'react';

import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';

import { BankAccount } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { FormSwitch } from '@ui/form/Switch/FormSwitch2';
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
          Bank transfer test
        </span>
        <span className='flex items-center'>
          <Tooltip label='Add new bank account' side='top' align='center'>
            <div>
              <AddAccountButton
                existingCurrencies={existingAccountCurrencies ?? []}
                organizationName={organizationName}
              />
            </div>
          </Tooltip>
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
