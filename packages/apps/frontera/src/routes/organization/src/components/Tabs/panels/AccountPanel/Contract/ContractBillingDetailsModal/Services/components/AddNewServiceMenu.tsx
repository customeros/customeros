import React from 'react';

import { observer } from 'mobx-react-lite';

import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { RefreshCcw02 } from '@ui/media/icons/RefreshCcw02';
import { CalendarDate } from '@ui/media/icons/CalendarDate';
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
          <MenuButton asChild>
            {isInline ? (
              <p>Add a service</p>
            ) : (
              <IconButton
                aria-label='Add a service'
                className='ml-1'
                size='xxs'
                variant='ghost'
                colorScheme='gray'
                icon={<Plus />}
              />
            )}
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
