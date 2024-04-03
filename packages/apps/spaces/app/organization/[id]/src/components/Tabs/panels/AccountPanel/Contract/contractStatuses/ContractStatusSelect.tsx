import React, { useMemo } from 'react';

import { UseMutationResult } from '@tanstack/react-query';
import {
  ContractEndModal,
  ContractStartModal,
} from 'app/organization/[id]/src/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Tag } from '@ui/presentation/Tag';
import { DotLive } from '@ui/media/icons/DotLive';
import { XSquare } from '@ui/media/icons/XSquare';
import { RefreshCw02 } from '@ui/media/icons/RefreshCw02';
import { SelectOption } from '@shared/types/SelectOptions';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';

import { contractOptionIcon } from './utils';

interface ContractStatusSelectProps {
  renewsAt?: string;
  contractId: string;
  status: ContractStatus;
  serviceStarted?: string;
  organizationName: string;
  nextInvoiceDate?: string;

  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

export const contractStatusOptions: SelectOption<ContractStatus>[] = [
  { label: 'Draft', value: ContractStatus.Draft },
  { label: 'Ended', value: ContractStatus.Ended },
  { label: 'Live', value: ContractStatus.Live },
  { label: 'Out of contract', value: ContractStatus.OutOfContract },
];

const statusColorScheme: Record<string, string> = {
  [ContractStatus.Live]: 'primary',
  [ContractStatus.Draft]: 'gray',
  [ContractStatus.Ended]: 'gray',
  [ContractStatus.OutOfContract]: 'warning',
};

export const ContractStatusSelect: React.FC<ContractStatusSelectProps> = ({
  status,
  renewsAt,
  contractId,
  organizationName,
  serviceStarted,
  onUpdateContract,
  nextInvoiceDate,
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

  const selected = contractStatusOptions.find((e) => e.value === status);
  const icon = contractOptionIcon?.[status];

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
          maxW={'auto'}
          color={`${statusColorScheme[status]}.800`}
          borderColor={`${statusColorScheme[status]}.800`}
          bg={`${statusColorScheme[status]}.50`}
        >
          <Tag
            as={Flex}
            alignItems='center'
            gap={1}
            variant='outline'
            whiteSpace='nowrap'
            colorScheme={statusColorScheme[status]}
            color={`${statusColorScheme[status]}.700`}
          >
            {icon && (
              <Flex alignItems='center' boxSize={3}>
                {icon}
              </Flex>
            )}

            {selected?.label}
          </Tag>
        </MenuButton>
        <MenuList minW={'150px'}>
          <MenuItem
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
