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
        <DotsVertical className='text-gray-400 size-4' />
      </MenuButton>
      <MenuList align='end' side='bottom' className='min-w-[150px]'>
        <MenuItem
          onClick={() => mutate({ id })}
          className='w-auto flex items-center'
        >
          <Archive className='mr-2 text-gray-500' />
          Archive account
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
