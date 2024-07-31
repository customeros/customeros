import { useState } from 'react';

import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useCreateBankAccountMutation } from '@settings/graphql/createBankAccount.generated';
import { currencyIcon } from '@settings/components/Tabs/panels/BillingPanel/components/utils';

import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { currencyOptions } from '@shared/util/currencyOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';

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
            icon={<Plus />}
            variant='ghost'
            colorScheme='gray'
            aria-label='Add account'
            onClick={() => setShowCurrencySelect(true)}
          />
        </Tooltip>
      )}

      {showCurrencySelect && (
        <Menu>
          <MenuButton>
            <IconButton
              size='md'
              icon={<Plus />}
              variant='ghost'
              colorScheme='gray'
              aria-label='Account currency'
            />
          </MenuButton>
          <MenuList>
            {currencyOptions.map((option) => (
              <MenuItem
                key={option.value}
                disabled={existingCurrencies?.indexOf(option.value) > -1}
                onSelect={() => {
                  mutate({
                    input: {
                      currency: option.value,
                      bankName: legalName,
                    },
                  });
                }}
              >
                {currencyIcon?.[option.value]}
                {option.value}
              </MenuItem>
            ))}
          </MenuList>
        </Menu>
      )}
    </>
  );
};
