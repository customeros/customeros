import React from 'react';

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
  organizationName: string;
  serviceStartedAt?: string;
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
];

export const ContractStatusSelect: React.FC<ContractStatusSelectProps> = ({
  status,
  renewsAt,
  contractId,
  organizationName,
  serviceStartedAt,
  onUpdateContract,
}) => {
  const { onOpen, onClose, isOpen } = useDisclosure({
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

  return (
    <>
      <Menu>
        <MenuButton
          disabled={status === ContractStatus.Draft}
          maxW={'auto'}
          color={status === ContractStatus.Live ? 'primary.800' : 'gray.800'}
          borderColor={
            status === ContractStatus.Live ? 'primary.800' : 'gray.800'
          }
          bg={status === ContractStatus.Live ? 'primary.50' : 'gray.50'}
        >
          <Tag
            as={Flex}
            alignItems='center'
            gap={1}
            variant='outline'
            colorScheme={
              selected?.value === ContractStatus.Live ? 'primary' : 'gray'
            }
            color={status === ContractStatus.Live ? 'primary.800' : 'gray.800'}
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
            onClick={status === ContractStatus.Live ? onOpen : onOpenStartModal}
          >
            {status === ContractStatus.Live ? (
              <>
                <XSquare color='gray.500' mr={1} />
                End contract...
              </>
            ) : (
              <>
                <DotLive color='gray.500' mr={1} />
                Make live
              </>
            )}
          </MenuItem>
        </MenuList>
      </Menu>

      <ContractEndModal
        isOpen={isOpen}
        onClose={onClose}
        contractId={contractId}
        organizationName={organizationName}
        renewsAt={renewsAt}
        serviceStartedAt={serviceStartedAt}
        onUpdateContract={onUpdateContract}
      />
      <ContractStartModal
        isOpen={isStartModalOpen}
        onClose={onCloseStartModal}
        contractId={contractId}
        organizationName={organizationName}
        serviceStartedAt={serviceStartedAt}
        onUpdateContract={onUpdateContract}
      />
    </>
  );
};
