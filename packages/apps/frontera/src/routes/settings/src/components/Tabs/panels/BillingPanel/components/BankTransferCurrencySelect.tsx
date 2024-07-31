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
          size='xs'
          variant='outline'
          colorScheme='gray'
          aria-label='Select currency'
          className='rounded-full size-6'
          icon={currencyIcon?.[currency || '']}
        />
      </MenuButton>
      <MenuList>
        {currencyOptions.map((option) => (
          <MenuItem
            id={id}
            key={option.value}
            disabled={existingCurrencies?.indexOf(option.value) > -1}
            onSelect={(e) => {
              onChange(e);
            }}
          >
            {currencyIcon?.[option.value]}
            {option.value}
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
  );
};
