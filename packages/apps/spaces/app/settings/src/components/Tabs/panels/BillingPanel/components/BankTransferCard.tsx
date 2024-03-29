'use client';

import React from 'react';
import { useForm } from 'react-inverted-form';

import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useUpdateBankAccountMutation } from '@settings/graphql/updateBankAccount.generated';
import { BankNameInput } from '@settings/components/Tabs/panels/BillingPanel/components/BankNameInput';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input/FormInput2';
import { FormMaskInput } from '@ui/form/Input/FormMaskInput';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { Currency, BankAccount, BankAccountUpdateInput } from '@graphql/types';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';

import { BankTransferMenu } from './BankTransferMenu';
import { BankTransferCurrencySelect } from './BankTransferCurrencySelect';

const bankOptions = {
  mask: '00 0000 0000 0000 0000 0000 0000 ',
  definitions: {
    '0': /[0-9]/,
  },
};

const sortCodeOptions = {
  mask: '00-00-00',
  definitions: {
    '0': /[0-9]/,
  },
};

export const BankTransferCard = ({
  account,
  existingCurrencies,
}: {
  account: BankAccount;
  existingCurrencies: Array<string>;
}) => {
  const formId = `bank-transfer-form-${account.metadata.id}`;
  const queryKey = useBankAccountsQuery.getKey();
  const queryClient = useQueryClient();

  const client = getGraphQLClient();
  const { mutate } = useUpdateBankAccountMutation(client, {
    onSuccess: () => {},
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });

  const updateBankAccountDebounced = useDebounce(
    (variables: Partial<BankAccountUpdateInput>) => {
      mutate({
        input: {
          id: account.metadata.id,
          ...variables,
        },
      });
    },
    500,
  );
  useForm<BankAccount>({
    formId,
    defaultValues: account,
    debug: true,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'bic':
          case 'sortCode':
          case 'iban':
          case 'routingNumber':
          case 'accountNumber':
          case 'bankName':
            updateBankAccountDebounced({
              [action.payload.name]: action.payload.value,
              currency: account.currency,
            });

            return next;
          case 'currency':
            mutate({
              input: {
                id: account.metadata.id,
                currency: action.payload.value?.value,
                sortCode: '',
                iban: '',
                routingNumber: '',
                accountNumber: '',
                bic: '',
              },
            });

            return next;

          default: {
            return next;
          }
        }
      }

      return next;
    },
  });

  return (
    <>
      <Card
        className='py-2 px-4 rounded-lg border-[1px] border-gray-200'

        // _hover={{
        //   '& #help-button': {
        //     visibility: 'visible',
        //   },
        // }}
      >
        <CardHeader className='p-0 pb-1 flex justify-between'>
          <BankNameInput formId={formId} metadata={account.metadata} />

          <div className='flex'>
            <BankTransferCurrencySelect
              existingCurrencies={existingCurrencies}
              currency={account.currency}
              formId={formId}
            />

            <BankTransferMenu id={account?.metadata?.id} />
          </div>
        </CardHeader>
        <CardContent className='p-0 gap-2'>
          {account.currency !== 'USD' && account.currency !== 'GBP' && (
            <>
              <FormMaskInput
                options={{ opts: bankOptions }}
                autoComplete='off'
                label='IBAN'
                placeholder='IBAN #'
                labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                name='iban'
                mb={1}
                formId={formId}
              />
              <FormInput
                autoComplete='off'
                label='BIC/Swift'
                placeholder='BIC/Swift'
                labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                name='bic'
                formId={formId}
              />
            </>
          )}
          <Flex pb={1} gap={2}>
            {account.currency === 'GBP' && (
              <>
                <FormMaskInput
                  options={{ opts: sortCodeOptions }}
                  autoComplete='off'
                  label='Sort code'
                  placeholder='Sort code'
                  labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                  name='sortCode'
                  formId={formId}
                  maxW='80px'
                />
                <FormMaskInput
                  options={{ opts: bankOptions }}
                  autoComplete='off'
                  label='Account number'
                  placeholder='Bank account #'
                  labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                  name='accountNumber'
                  formId={formId}
                />
              </>
            )}
          </Flex>
          {account.currency === 'USD' && (
            <>
              <FormInput
                autoComplete='off'
                className='mb-1'
                label='Routing number'
                placeholder='Routing number'
                labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                name='routingNumber'
                formId={formId}
              />
              <FormMaskInput
                options={{ opts: bankOptions }}
                autoComplete='off'
                label='Account number'
                placeholder='Bank account #'
                labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                name='accountNumber'
                formId={formId}
              />
            </>
          )}
          {account.allowInternational &&
            (account.currency === 'USD' || account.currency === 'GBP') && (
              <FormInput
                autoComplete='off'
                label='BIC/Swift'
                placeholder='BIC/Swift'
                labelProps={{ className: 'text-sm mb-0 font-semibold' }}
                name='bic'
                formId={formId}
              />
            )}

          {(account.allowInternational ||
            ![Currency.Gbp, Currency.Usd, Currency.Eur].includes(
              account?.currency as Currency,
            )) && (
            <FormAutoresizeTextarea
              autoComplete='off'
              label='Other details'
              placeholder='Other details'
              name='otherDetails'
              formId={formId}
            />
          )}
        </CardContent>
      </Card>
    </>
  );
};
