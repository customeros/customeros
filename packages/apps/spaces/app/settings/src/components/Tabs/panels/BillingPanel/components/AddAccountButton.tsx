'use client';

import React, { useMemo, useState } from 'react';

import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useCreateBankAccountMutation } from '@settings/graphql/createBankAccount.generated';
import {
  currencyIcon,
  mapCurrencyToOptions,
} from '@settings/components/Tabs/panels/BillingPanel/components/utils';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Select } from '@ui/form/SyncSelect';
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const AddAccountButton = ({
  existingCurrencies,
}: {
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
  const currencyOptions = useMemo(() => mapCurrencyToOptions(), []);

  return (
    <>
      {!showCurrencySelect && (
        <IconButton
          size='xs'
          color='gray.400'
          icon={<Plus />}
          variant='ghost'
          aria-label='Add account'
          onClick={() => setShowCurrencySelect(true)}
        />
      )}

      {showCurrencySelect && (
        <Select
          placeholder='Account currency'
          name='renewalCycle'
          blurInputOnSelect
          onChange={(e) => {
            mutate({
              input: {
                currency: e.value,
              },
            });
          }}
          onBlur={() => setShowCurrencySelect(false)}
          options={[
            {
              options: [
                { label: 'USD', value: 'USD' },
                { label: 'GBP', value: 'GBP' },
                { label: 'EUR', value: 'EUR' },
              ],
            },

            ...currencyOptions,
          ]}
          isOptionDisabled={(option, selectValue) =>
            existingCurrencies?.indexOf(option.value) > -1
          }
          getOptionLabel={(option) => {
            const alreadyExists =
              existingCurrencies?.indexOf(option.value) > -1;

            return (
              <Tooltip
                label={alreadyExists ? 'Already used on another account' : ''}
              >
                <Flex alignItems='center'>
                  {currencyIcon?.[option.value]}

                  <Text className='option-label'>{option.label}</Text>
                </Flex>
              </Tooltip>
            ) as unknown as string;
          }}
          chakraStyles={{
            container: (props, state) => {
              return {
                ...props,
                w: '100%',
                overflow: 'visible',
                willChange: 'width',
                transition: 'width 0.2s',
                _hover: { cursor: 'pointer' },
                minHeight: '10px',
                maxH: '24px',
                height: '24px',
              };
            },
            control: (props, state) => {
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
