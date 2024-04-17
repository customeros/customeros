'use client';

import React from 'react';

import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useDeleteBankAccountMutation } from '@settings/graphql/deleteBankAccount.generated';

import { Archive } from '@ui/media/icons/Archive';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const BankTransferMenu = ({ id }: { id: string }) => {
  const queryKey = useBankAccountsQuery.getKey();
  const queryClient = useQueryClient();

  const client = getGraphQLClient();
  const { mutate } = useDeleteBankAccountMutation(client, {
    onSuccess: () => {},
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });

  return (
    <Menu>
      <MenuButton className='mb-1 ml-1'>
        <DotsVertical color='gray.400' boxSize={4} />
      </MenuButton>
      <MenuList side='bottom' align='end' className='min-w-[150px]'>
        <MenuItem
          className='w-auto flex items-center'
          onClick={() => mutate({ id })}
        >
          <Archive mr={2} color='gray.500' />
          Archive account
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
