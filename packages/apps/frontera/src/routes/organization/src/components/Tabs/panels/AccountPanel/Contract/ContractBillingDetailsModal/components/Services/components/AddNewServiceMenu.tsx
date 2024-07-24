import React from 'react';

import { observer } from 'mobx-react-lite';

import { BilledType } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { RefreshCcw02 } from '@ui/media/icons/RefreshCcw02.tsx';
import { CalendarDate } from '@ui/media/icons/CalendarDate.tsx';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

interface AddNewServiceMenuProps {
  isInline?: boolean;
  contractId: string;
}
//TODO K
export const AddNewServiceMenu: React.FC<AddNewServiceMenuProps> = observer(
  ({ isInline, contractId }) => {
    const store = useStore();
    const contractLineItemsStore = store.contractLineItems;

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
                variant='outline'
                colorScheme='gray'
                icon={<Plus />}
              />
            )}
          </MenuButton>

          <MenuList align='end' side='bottom' className='p-0'>
            <MenuItem
              onClick={() =>
                contractLineItemsStore.create({
                  billingCycle: BilledType.Monthly,
                  contractId,
                  description: 'Unnamed',
                })
              }
              className='flex items-center text-base'
            >
              <RefreshCcw02 className='mr-2 text-gray-500' />
              Subscription
            </MenuItem>
            <MenuItem
              onClick={() =>
                contractLineItemsStore.create({
                  billingCycle: BilledType.Once,
                  contractId,
                  description: 'Unnamed',
                })
              }
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
