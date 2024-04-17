import React, { ReactNode } from 'react';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { ContractStatus } from '@graphql/types';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface ContractStatusSelectProps {
  status: ContractStatus;
  statusContent: ReactNode;
  onOpenEditModal: () => void;
  onHandleStatusChange: () => void;
}

export const ContractMenu: React.FC<ContractStatusSelectProps> = ({
  status,
  onOpenEditModal,
  onHandleStatusChange,
  statusContent,
}) => {
  return (
    <>
      <Menu>
        <MenuButton
          className={cn(
            `flex items-center max-h-5 p-1 hover:bg-gray-100 rounded`,
          )}
        >
          <DotsVertical color='gray.400' />
        </MenuButton>
        <MenuList align='end' side='bottom'>
          <MenuItem onClick={onOpenEditModal} className='flex items-center'>
            <Edit03 mr={2} color='gray.500' />
            Edit contract
          </MenuItem>
          {status !== ContractStatus.Scheduled && (
            <MenuItem
              className='flex items-center'
              onClick={onHandleStatusChange}
            >
              {statusContent}
            </MenuItem>
          )}
        </MenuList>
      </Menu>
    </>
  );
};
