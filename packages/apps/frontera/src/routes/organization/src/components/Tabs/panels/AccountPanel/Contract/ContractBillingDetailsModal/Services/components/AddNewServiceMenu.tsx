import React from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { RefreshCcw02 } from '@ui/media/icons/RefreshCcw02';
import { CalendarDate } from '@ui/media/icons/CalendarDate';
import { ghostButton } from '@ui/form/Button/Button.variants';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { useEditContractModalStores } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/EditContractModalStores';

interface AddNewServiceMenuProps {
  isInline?: boolean;
}

export const AddNewServiceMenu: React.FC<AddNewServiceMenuProps> = observer(
  ({ isInline }) => {
    const { serviceFormStore } = useEditContractModalStores();

    return (
      <>
        <Menu>
          <MenuButton
            className={cn(
              ghostButton({ colorScheme: 'gray' }),
              `flex items-center max-h-5 p-1 ml-[5px] hover:bg-gray-100 rounded`,
            )}
          >
            {isInline ? <p>Add a service</p> : <Plus className='size-3' />}
          </MenuButton>
          <MenuList align='end' side='bottom' className='p-0'>
            <MenuItem
              onClick={() => serviceFormStore.addService(null, true)}
              className='flex items-center text-base'
            >
              <RefreshCcw02 className='mr-2 text-gray-500' />
              Subscription
            </MenuItem>
            <MenuItem
              onClick={() => serviceFormStore.addService(null)}
              className='flex items-center text-base'
            >
              <CalendarDate className='mr-2 text-gray-500' />
              One-time
            </MenuItem>
          </MenuList>
        </Menu>
      </>
    );
  },
);
