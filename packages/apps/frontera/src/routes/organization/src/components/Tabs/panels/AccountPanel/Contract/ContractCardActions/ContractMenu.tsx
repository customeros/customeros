import React, { ReactNode } from 'react';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { ContractStatus } from '@graphql/types';
import { Trash01 } from '@ui/media/icons/Trash01.tsx';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Divider } from '@ui/presentation/Divider/Divider';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

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
          data-test='contract-menu-dots'
          className={cn(
            `flex items-center max-h-5 p-1 hover:bg-gray-100 rounded`,
          )}
        >
          <DotsVertical className='text-gray-400' />
        </MenuButton>
        <MenuList align='end' side='bottom' className='p-0'>
          <MenuItem
            onClick={onOpenEditModal}
            className='flex items-center text-base'
            data-test='contract-menu-edit-contract'
          >
            <Edit03 className='mr-1 text-gray-500' />
            Edit contract
          </MenuItem>

          {status !== ContractStatus.Scheduled && (
            <>
              {status === ContractStatus.Live && (
                <MenuItem
                  className='flex items-center text-base'
                  onClick={() =>
                    onStatusModalOpen(ContractStatusModalMode.Renew)
                  }
                >
                  <RefreshCw05 className='text-gray-500 mr-1' />
                  Renew contract
                </MenuItem>
              )}
              <Divider className='my-0.5' />
              <MenuItem
                onClick={onHandleStatusChange}
                className='flex items-center text-base'
              >
                {statusContent}
              </MenuItem>
            </>
          )}
          {status == ContractStatus.Draft && (
            <MenuItem
              className='flex items-center text-base'
              data-test='contract-menu-delete-contract'
              onClick={() => onStatusModalOpen(ContractStatusModalMode.Delete)}
            >
              <Trash01 className='mr-1 text-gray-500' />
              Delete contract
            </MenuItem>
          )}
        </MenuList>
      </Menu>
    </>
  );
};
