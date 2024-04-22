import React, { ReactNode } from 'react';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { ContractStatus } from '@graphql/types';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Divider } from '@ui/presentation/Divider/Divider';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

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
  const { onStatusModalOpen } = useContractModalStatusContext();

  return (
    <>
      <Menu>
        <MenuButton
          className={cn(
            `flex items-center max-h-5 p-1 hover:bg-gray-100 rounded`,
          )}
        >
          <DotsVertical className='text-gray-400' />
        </MenuButton>
        <MenuList align='end' side='bottom' className='p-0'>
          <MenuItem onClick={onOpenEditModal} className='flex items-center'>
            <Edit03 className='mr-2 text-gray-500' />
            Edit contract
          </MenuItem>

          {status === ContractStatus.Live && (
            <MenuItem
              className='flex items-center'
              onClick={() => onStatusModalOpen(ContractStatusModalMode.Renew)}
            >
              <RefreshCw05 color='gray.500' mr={2} />
              Renew contract
            </MenuItem>
          )}
          <Divider className='my-0.5' />
          {status !== ContractStatus.Scheduled && (
            <MenuItem
              className='flex items-center '
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
