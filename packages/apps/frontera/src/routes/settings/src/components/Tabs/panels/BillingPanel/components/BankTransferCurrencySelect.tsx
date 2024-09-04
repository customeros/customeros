import { IconButton } from '@ui/form/IconButton/IconButton';
import { currencyOptions } from '@shared/util/currencyOptions';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

import { currencyIcon } from './utils';

export const BankTransferCurrencySelect = ({
  currency,
  existingCurrencies,
  onChange,
}: {
  currency?: string | null;
  existingCurrencies: Array<string>;
  onChange: (value: string) => void;
}) => {
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
            key={option.value}
            disabled={existingCurrencies?.indexOf(option.value) > -1}
            onSelect={() => {
              onChange(option.value);
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
