'use client';

import React from 'react';
import { useForm } from 'react-inverted-form';

import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useUpdateBankAccountMutation } from '@settings/graphql/updateBankAccount.generated';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input/FormInput2';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { FormMaskInput } from '@ui/form/Input/FormMaskInput';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { Currency, BankAccount, BankAccountUpdateInput } from '@graphql/types';

import { BankTransferMenu } from './BankTransferMenu';
import { BankTransferCurrencySelect } from './BankTransferCurrencySelect';

const sortCode = {
  mask: '00-00-00',
  definitions: {
    '0': /[0-9]/,
  },
};

const accountNumber = {
  mask: '00 0000 0000 0000 0000 0000 0000 ',
  definitions: {
    '0': /[0-9]/,
  },
};

export const BankTransferCard = ({ account }: { account: BankAccount }) => {
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
        // _hover={{
        //   '& #help-button': {
        //     visibility: 'visible',
        //   },
        // }}
        className={'py-2 px-4 border-b border-gray-200 rounded-lg bg-white '}
      >
        <CardHeader className='p-0 grid grid-cols-3'>
          <div className='col-span-2'>
            <FormInput
              formId={formId}
              name='bankName'
              label='Bank Name'
              labelProps={{
                style: { display: 'none' },
              }}
              variant='unstyled'
              autoComplete='off'
              placeholder='Bank name'
              className='font-semibold '
            />
          </div>

          <div className='flex justify-end'>
            <BankTransferCurrencySelect
              currency={account.currency}
              formId={formId}
            />

            <BankTransferMenu
              id={account?.metadata?.id}
              allowInternational={account.allowInternational}
              currency={account?.currency}
            />
          </div>
        </CardHeader>
        <CardContent className='p-0 gap-2'>
          <Flex pb={1} gap={2}>
            {account.currency === 'GBP' && (
              <>
                <FormMaskInput
                  options={{ opts: sortCode }}
                  autoComplete='off'
                  label='Sort code'
                  placeholder='Sort code'
                  labelProps={{
                    className: 'font-semibold mb-0 text-sm',
                  }}
                  name='sortCode'
                  formId={formId}
                  className='max-w-20'
                />
                <FormMaskInput
                  options={{ opts: accountNumber }}
                  autoComplete='off'
                  label='Account number'
                  placeholder='Bank account #'
                  labelProps={{
                    className: 'font-semibold mb-0 text-sm',
                  }}
                  name='accountNumber'
                  formId={formId}
                />
              </>
            )}
            {account.currency !== 'USD' && account.currency !== 'GBP' && (
              <>
                <FormInput
                  autoComplete='off'
                  label='BIC/Swift'
                  placeholder='BIC/Swift'
                  labelProps={{
                    className: 'font-semibold mb-0 text-sm',
                  }}
                  name='bic'
                  formId={formId}
                />

                <FormMaskInput
                  options={{ opts: accountNumber }}
                  autoComplete='off'
                  label='Iban'
                  placeholder='Iban #'
                  labelProps={{
                    className: 'font-semibold mb-0 text-sm',
                  }}
                  name='iban'
                  formId={formId}
                />
              </>
            )}
          </Flex>
          {account.currency === 'USD' && (
            <>
              <FormInput
                autoComplete='off'
                label='Routing number'
                placeholder='Routing number'
                labelProps={{
                  className: 'font-semibold mb-0 text-sm',
                }}
                name='routingNumber'
                formId={formId}
              />

              <FormMaskInput
                options={{
                  opts: accountNumber,
                }}
                autoComplete='off'
                label='Account number'
                placeholder='Bank account #'
                labelProps={{
                  className: 'font-semibold mb-0 text-sm',
                }}
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
                labelProps={{
                  className: 'font-semibold mb-0 text-sm',
                }}
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
