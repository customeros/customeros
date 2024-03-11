'use client';

import React, { useMemo } from 'react';
import { useField } from 'react-inverted-form';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Select } from '@ui/form/SyncSelect';

import { currencyIcon, mapCurrencyToOptions } from './utils';

export const BankTransferCurrencySelect = ({
  formId,
  currency,
}: {
  formId: string;
  currency?: string | null;
}) => {
  const { getInputProps } = useField('currency', formId);
  const { id, onChange } = getInputProps();
  const currencyOptions = useMemo(() => mapCurrencyToOptions(), []);

  return (
    <Select
      id={id}
      placeholder='Account currency'
      name='renewalCycle'
      blurInputOnSelect
      onChange={(e) => {
        onChange(e);
      }}
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
      getOptionLabel={(option) => {
        return (
          <Flex alignItems='center'>
            {currencyIcon?.[option.value]}
            <Text className='option-label'>{option.value}</Text>
          </Flex>
        ) as unknown as string;
      }}
      defaultValue={{ label: currency, value: currency }}
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
              width: '24px',
              maxW: '24px',
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
              width: '24px',
              maxW: '24px',
              height: '24px',
              maxH: '24px',
              borderRadius: '50%',
              border: '1px solid',
              borderColor: 'gray.200',

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
  );
};
