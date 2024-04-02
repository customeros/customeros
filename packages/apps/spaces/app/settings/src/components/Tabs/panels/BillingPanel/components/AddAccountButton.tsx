'use client';

import React, { useState } from 'react';

import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useCreateBankAccountMutation } from '@settings/graphql/createBankAccount.generated';
import { currencyIcon } from '@settings/components/Tabs/panels/BillingPanel/components/utils';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Select } from '@ui/form/SyncSelect';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { currencyOptions } from '@shared/util/currencyOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const AddAccountButton = ({
  existingCurrencies,
  legalName,
}: {
  legalName?: string | null;
  existingCurrencies: Array<string>;
}) => {
  const [showCurrencySelect, setShowCurrencySelect] = useState(false);
  const queryKey = useBankAccountsQuery.getKey();
  const queryClient = useQueryClient();
  const client = getGraphQLClient();

  const { mutate } = useCreateBankAccountMutation(client, {
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey });
    },
    onSettled: () => {
      setShowCurrencySelect(false);
    },
  });

  return (
    <>
      {!showCurrencySelect && (
        <Tooltip label='Add new bank account'>
          <IconButton
            size='xs'
            color='gray.400'
            icon={<Plus />}
            variant='ghost'
            aria-label='Add account'
            onClick={() => setShowCurrencySelect(true)}
          />
        </Tooltip>
      )}

      {showCurrencySelect && (
        <Select
          placeholder='Account currency'
          name='bankAccountCurrency'
          defaultMenuIsOpen
          menuIsOpen
          blurInputOnSelect
          onChange={(e) => {
            mutate({
              input: {
                currency: e.value,
                bankName: legalName,
              },
            });
          }}
          onBlur={() => setShowCurrencySelect(false)}
          options={currencyOptions}
          isOptionDisabled={(option, selectValue) =>
            existingCurrencies?.indexOf(option.value) > -1
          }
          formatOptionLabel={(option, { context }) => {
            return (
              <Flex alignItems='center'>
                <Flex
                  justifyContent={context === 'value' ? 'center' : 'flex-end'}
                  alignItems='center'
                >
                  {currencyIcon?.[option.value]}
                </Flex>
                <Text className='option-label' ml={3}>
                  {option.value}
                </Text>
              </Flex>
            );
          }}
          chakraStyles={{
            container: (props, state) => {
              if (
                !state?.selectProps?.menuIsOpen &&
                state.hasValue &&
                !state.isFocused
              ) {
                return {
                  display: 'flex',
                  alignItems: 'center',
                  width: 'fit-content',
                  maxW: 'fit-content',
                  willChange: 'width',
                  transition: 'width 0.2s',
                };
              }

              return {
                ...props,
                w: '100%',
                overflow: 'visible',
                willChange: 'width',
                transition: 'width 0.2s',
                _hover: { cursor: 'pointer' },
              };
            },
            control: (props, state) => {
              if (
                !state?.selectProps?.menuIsOpen &&
                state.hasValue &&
                !state.isFocused
              ) {
                return {
                  height: '24px',
                  maxH: '24px',
                  width: 'max-content',
                  minW: '24px',
                  borderRadius: '30px',
                  border: '1px solid',
                  borderColor: 'gray.200',
                  padding: '2px',

                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center',
                  fontSize: '12px',

                  '& .option-label': {
                    display: 'none',
                  },
                  '& svg': {
                    marginLeft: '1px',
                    height: '12px',
                  },
                };
              }

              return {
                ...props,
                w: '100%',
                overflow: 'visible',
                _hover: { cursor: 'pointer' },
              };
            },
            groupHeading: (props) => ({
              display: 'none',
            }),
            group: (props) => ({
              borderBottom: '1px solid',
              borderColor: 'gray.200',
            }),
          }}
        />
      )}
    </>
  );
};
