import React from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { XSquare } from '@ui/media/icons/XSquare';
import { BracketsPlus } from '@ui/media/icons/BracketsPlus';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { useEditContractModalStores } from '../../stores/EditContractModalStores';

interface ServiceItemMenuProps {
  id: string;
  closed?: boolean;
  allowAddModification?: boolean;
  type: 'subscription' | 'one-time';
  handleCloseService: (isClosed: boolean) => void;
}

export const ServiceItemMenu: React.FC<ServiceItemMenuProps> = observer(
  ({ id, type, allowAddModification, handleCloseService }) => {
    const { serviceFormStore } = useEditContractModalStores();

    return (
      <>
        <Menu>
          <MenuButton
            className={cn(
              `flex items-center max-h-5 p-1 py-2 hover:bg-gray-100 rounded translate-x-2`,
            )}
          >
            <DotsVertical className='text-gray-400' />
          </MenuButton>
          <MenuList align='end' side='bottom' className='p-0'>
            {allowAddModification && (
              <MenuItem
                onClick={() =>
                  serviceFormStore.addService(id, type === 'subscription')
                }
                className='flex items-center text-base'
              >
                <BracketsPlus className='mr-2 text-gray-500' />
                Add modification
              </MenuItem>
            )}

            <MenuItem
              onClick={() => handleCloseService(true)}
              className='flex items-center text-base'
            >
              <XSquare className='mr-2 text-gray-500' />
              End the service
            </MenuItem>
          </MenuList>
        </Menu>
      </>
    );
  },
);