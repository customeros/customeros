'use client';

import React from 'react';

import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useUpdateBankAccountMutation } from '@settings/graphql/updateBankAccount.generated';
import { useDeleteBankAccountMutation } from '@settings/graphql/deleteBankAccount.generated';

import { Currency } from '@graphql/types';
import { Globe04 } from '@ui/media/icons/Globe04';
import { Archive } from '@ui/media/icons/Archive';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const BankTransferMenu = ({
  id,
  allowInternational,
  currency,
}: {
  id: string;
  currency?: Currency | null;
  allowInternational: boolean;
}) => {
  const queryKey = useBankAccountsQuery.getKey();
  const queryClient = useQueryClient();

  const client = getGraphQLClient();
  const { mutate } = useDeleteBankAccountMutation(client, {
    onSuccess: () => {},
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });
  const { mutate: update } = useUpdateBankAccountMutation(client, {
    onSuccess: () => {},
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });
  const showInternational = [Currency.Usd, Currency.Gbp, Currency.Eur].includes(
    currency as Currency,
  );
  const toggleAllowInternational = () => {
    update({
      input: {
        id,
        allowInternational: !allowInternational,
        bic: '',
        otherDetails: '',
      },
    });
  };

  return (
    <Menu>
      <MenuButton className='mb-1 ml-1'>
        <DotsVertical color='gray.400' boxSize={4} />
      </MenuButton>
      <MenuList className='min-w-[150px]'>
        {!allowInternational && showInternational && (
          <MenuItem onClick={toggleAllowInternational}>
            <Globe04 mr={2} color='gray.500' />
            Add international
          </MenuItem>
        )}
        {allowInternational && showInternational && (
          <MenuItem onClick={toggleAllowInternational}>
            <MinusCircle mr={2} color='gray.500' />
            Remove international
          </MenuItem>
        )}

        <MenuItem onClick={() => mutate({ id })}>
          <Archive mr={2} color='gray.500' />
          Archive account
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
