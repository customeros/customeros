import { observer } from 'mobx-react-lite';
import { currencyIcon } from '@settings/components/Tabs/panels/BillingPanel/components/utils';

import { Plus } from '@ui/media/icons/Plus';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { currencyOptions } from '@shared/util/currencyOptions';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

export const AddAccountButton = observer(
  ({
    existingCurrencies,
    legalName,
  }: {
    legalName?: string | null;
    existingCurrencies: Array<string>;
  }) => {
    const store = useStore();

    return (
      <>
        <Menu>
          <MenuButton>
            <Tooltip label='Add new bank account'>
              <IconButton
                size='sm'
                variant='ghost'
                colorScheme='gray'
                aria-label='Account currency'
                icon={<Plus className='size-4' />}
              />
            </Tooltip>
          </MenuButton>
          <MenuList>
            {currencyOptions.map((option) => (
              <MenuItem
                key={option.value}
                disabled={existingCurrencies?.indexOf(option.value) > -1}
                onSelect={() => {
                  store.settings.bankAccounts.create({
                    currency: option.value,
                    bankName: `${legalName} account`,
                  });
                }}
              >
                {currencyIcon?.[option.value]}
                {option.value}
              </MenuItem>
            ))}
          </MenuList>
        </Menu>
      </>
    );
  },
);
