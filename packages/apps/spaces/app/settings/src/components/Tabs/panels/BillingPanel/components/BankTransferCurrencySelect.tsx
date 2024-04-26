'use client';

import React from 'react';
import { useField } from 'react-inverted-form';

import { IconButton } from '@ui/form/IconButton/IconButton';
import { currencyOptions } from '@shared/util/currencyOptions';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

import { currencyIcon } from './utils';

export const BankTransferCurrencySelect = ({
  formId,
  currency,
  existingCurrencies,
}: {
  formId: string;
  currency?: string | null;
  existingCurrencies: Array<string>;
}) => {
  const { getInputProps } = useField('currency', formId);
  const { id, onChange } = getInputProps();

  return (
    <Menu>
      <MenuButton>
        <IconButton
          aria-label='Select currency'
          variant='outline'
          colorScheme='gray'
          icon={currencyIcon?.[currency || '']}
          size='xs'
          className='rounded-full size-6'
        />
      </MenuButton>
      <MenuList>
        {currencyOptions.map((option) => (
          <MenuItem
            id={id}
            key={option.value}
            onSelect={(e) => {
              onChange(e);
            }}
            disabled={existingCurrencies?.indexOf(option.value) > -1}
          >
            {currencyIcon?.[option.value]}
            {option.value}
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
  );
};
