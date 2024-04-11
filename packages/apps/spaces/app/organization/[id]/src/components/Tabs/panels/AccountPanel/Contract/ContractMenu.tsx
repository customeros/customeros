import React, { useMemo } from 'react';

import { UseMutationResult } from '@tanstack/react-query';
import {
  ContractEndModal,
  ContractStartModal,
} from 'app/organization/[id]/src/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals';

import { useDisclosure } from '@ui/utils';
import { Edit03 } from '@ui/media/icons/Edit03';
import { DotLive } from '@ui/media/icons/DotLive';
import { XSquare } from '@ui/media/icons/XSquare';
import { RefreshCw02 } from '@ui/media/icons/RefreshCw02';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';

interface ContractStatusSelectProps {
  renewsAt?: string;
  contractId: string;
  status: ContractStatus;
  serviceStarted?: string;
  organizationName: string;
  nextInvoiceDate?: string;
  onOpenEditModal: () => void;

  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

const statusColorScheme: Record<string, string> = {
  [ContractStatus.Live]: 'primary',
  [ContractStatus.Draft]: 'gray',
  [ContractStatus.Ended]: 'gray',
  [ContractStatus.OutOfContract]: 'warning',
};

export const ContractMenu: React.FC<ContractStatusSelectProps> = ({
  status,
  renewsAt,
  contractId,
  organizationName,
  serviceStarted,
  onUpdateContract,
  nextInvoiceDate,
  onOpenEditModal,
}) => {
  const {
    onOpen: onOpenEndModal,
    onClose,
    isOpen,
  } = useDisclosure({
    id: 'end-contract-modal',
  });
  const {
    onOpen: onOpenStartModal,
    onClose: onCloseStartModal,
    isOpen: isStartModalOpen,
  } = useDisclosure({
    id: 'start-contract-modal',
  });

  const getStatusDisplay = useMemo(() => {
    let icon, text;
    switch (status) {
      case ContractStatus.Live:
        icon = <XSquare color='gray.500' mr={1} />;
        text = 'End contract...';
        break;
      case ContractStatus.Draft:
      case ContractStatus.Ended:
        icon = <DotLive color='gray.500' mr={1} />;
        text = 'Make live';
        break;
      case ContractStatus.OutOfContract:
        icon = <RefreshCw02 color='gray.500' mr={2} />;
        text = 'Renew contract';
        break;
      default:
        icon = null;
        text = null;
    }

    return (
      <>
        {icon}
        {text}
      </>
    );
  }, [status]);

  return (
    <>
      <Menu>
        <MenuButton
          className={`flex items-center max-h-5 text-${statusColorScheme[status]}.800} border-${statusColorScheme[status]}.800} bg-${statusColorScheme[status]}.50}`}
        >
          <DotsVertical color='gray.400' />
        </MenuButton>
        <MenuList align='end' side='bottom'>
          <MenuItem onClick={onOpenEditModal} className='flex items-center'>
            <Edit03 mr={2} color='gray.500' />
            Edit contract
          </MenuItem>
          <MenuItem
            className='flex items-center'
            onClick={
              status === ContractStatus.Live ? onOpenEndModal : onOpenStartModal
            }
          >
            {getStatusDisplay}
          </MenuItem>
        </MenuList>
      </Menu>

      <ContractEndModal
        isOpen={isOpen}
        onClose={onClose}
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
