import React, { useState } from 'react';

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
            colorScheme='gray'
            icon={<Plus />}
            variant='ghost'
            aria-label='Add account'
            onClick={() => setShowCurrencySelect(true)}
          />
        </Tooltip>
      )}

      {showCurrencySelect && (
        <Menu>
          <MenuButton>
            <IconButton
              icon={<Plus />}
              aria-label='Account currency'
              variant='ghost'
              colorScheme='gray'
              size='md'
            />
          </MenuButton>
          <MenuList>
            {currencyOptions.map((option) => (
              <MenuItem
                key={option.value}
                onSelect={() => {
                  mutate({
                    input: {
                      currency: option.value,
                      bankName: legalName,
                    },
                  });
                }}
                disabled={existingCurrencies?.indexOf(option.value) > -1}
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
