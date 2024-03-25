'use client';

import React from 'react';
import { useForm } from 'react-inverted-form';

import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useUpdateBankAccountMutation } from '@settings/graphql/updateBankAccount.generated';
import { BankNameInput } from '@settings/components/Tabs/panels/BillingPanel/components/BankNameInput';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';
import { Currency, BankAccount, BankAccountUpdateInput } from '@graphql/types';

import { BankTransferMenu } from './BankTransferMenu';
import { SortCodeInput, BankAccountInput } from './inputs';
import { BankTransferCurrencySelect } from './BankTransferCurrencySelect';

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
        py={2}
        px={4}
        borderRadius='lg'
        boxShadow='none'
        border='1px solid'
        borderColor='gray.200'
        _hover={{
          '& #help-button': {
            visibility: 'visible',
          },
        }}
      >
        <CardHeader p='0' pb={1} as={Flex}>
          <BankNameInput formId={formId} metadata={account.metadata} />
          <BankTransferCurrencySelect
            currency={account.currency}
            formId={formId}
          />

          <BankTransferMenu id={account?.metadata?.id} />
        </CardHeader>
        <CardBody p={0} gap={2}>
          <Flex pb={1} gap={2}>
            {account.currency === 'GBP' && (
              <>
                <SortCodeInput
                  autoComplete='off'
                  label='Sort code'
                  placeholder='Sort code'
                  isLabelVisible
                  labelProps={{
                    fontSize: 'sm',
                    mb: 0,
                    fontWeight: 'semibold',
                  }}
                  name='sortCode'
                  formId={formId}
                  maxW='80px'
                />
                <BankAccountInput
                  autoComplete='off'
                  label='Account number'
                  placeholder='Bank account #'
                  isLabelVisible
                  labelProps={{
                    fontSize: 'sm',
                    mb: 0,
                    fontWeight: 'semibold',
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
                  isLabelVisible
                  labelProps={{
                    fontSize: 'sm',
                    mb: 0,
                    fontWeight: 'semibold',
                  }}
                  name='bic'
                  formId={formId}
                />
                <BankAccountInput
                  autoComplete='off'
                  label='Iban'
                  placeholder='Iban #'
                  isLabelVisible
                  labelProps={{
                    fontSize: 'sm',
                    mb: 0,
                    fontWeight: 'semibold',
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
                isLabelVisible
                labelProps={{
                  fontSize: 'sm',
                  mb: 0,
                  fontWeight: 'semibold',
                }}
                name='routingNumber'
                formId={formId}
              />
              <BankAccountInput
                autoComplete='off'
                label='Account number'
                placeholder='Bank account #'
                isLabelVisible
                labelProps={{
                  fontSize: 'sm',
                  mb: 0,
                  fontWeight: 'semibold',
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
                isLabelVisible
                labelProps={{
                  fontSize: 'sm',
                  mb: 0,
                  fontWeight: 'semibold',
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
        </CardBody>
      </Card>
    </>
  );
};
