'use client';

import React from 'react';

import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { BankAccount } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { AddAccountButton } from './AddAccountButton';
import { BankTransferCard } from './BankTransferCard';

export const BankTransferAccountList = () => {
  const client = getGraphQLClient();
  const { data } = useBankAccountsQuery(client);

  const existingAccountCurrencies = data?.bankAccounts.map(
    (account) => account.currency as string,
  );

  return (
    <>
      <Flex justifyContent='space-between' alignItems='center'>
        <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>
          Bank transfer
        </Text>
        <Flex alignItems='center'>
          <AddAccountButton
            existingCurrencies={existingAccountCurrencies ?? []}
          />
          <Box>
            {/*<FormSwitch name='canPayWithCard' formId={formId} size='sm' />*/}
          </Box>
        </Flex>
      </Flex>

      {data?.bankAccounts?.map((account) => (
        <React.Fragment key={account.metadata.id}>
          <BankTransferCard account={account as BankAccount} />
        </React.Fragment>
      ))}
    </>
  );
};
