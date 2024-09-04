import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const BankTransferMenu = observer(({ id }: { id: string }) => {
  const store = useStore();

  return (
    <Menu>
      <MenuButton className='mb-1 ml-1'>
        <DotsVertical className='text-gray-400 size-4' />
      </MenuButton>
      <MenuList align='end' side='bottom' className='min-w-[150px]'>
        <MenuItem
          className='w-auto flex items-center'
          onClick={() => store.settings.bankAccounts.remove(id)}
        >
          <Archive className='mr-2 text-gray-500' />
          Archive account
        </MenuItem>
      </MenuList>
    </Menu>
  );
});
