import React, { ReactNode } from 'react';

import { UseMutationResult } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  ContractEndModal,
  ContractStartModal,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals';

interface ContractStatusSelectProps {
  renewsAt?: string;
  contractId: string;
  status: ContractStatus;
  serviceStarted?: string;
  isEndModalOpen: boolean;
  organizationName: string;
  nextInvoiceDate?: string;
  statusContent: ReactNode;
  isStartModalOpen: boolean;
  onOpenEditModal: () => void;

  onCloseEndModal: () => void;
  onCloseStartModal: () => void;
  onHandleStatusChange: () => void;
  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

export const ContractMenu: React.FC<ContractStatusSelectProps> = ({
  status,
  renewsAt,
  contractId,
  organizationName,
  serviceStarted,
  onUpdateContract,
  nextInvoiceDate,
  onOpenEditModal,
  onHandleStatusChange,
  isStartModalOpen,
  isEndModalOpen,
  statusContent,
  onCloseEndModal,
  onCloseStartModal,
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

      <ContractEndModal
        isOpen={isEndModalOpen}
        onClose={onCloseEndModal}
        contractId={contractId}
        organizationName={organizationName}
        renewsAt={renewsAt}
        serviceStarted={serviceStarted}
        onUpdateContract={onUpdateContract}
        nextInvoiceDate={nextInvoiceDate}
      />
      <ContractStartModal
        isOpen={isStartModalOpen}
        onClose={onCloseStartModal}
        contractId={contractId}
        organizationName={organizationName}
        serviceStarted={serviceStarted}
        onUpdateContract={onUpdateContract}
      />
    </>
  );
};
