import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Select } from '@ui/form/SyncSelect';
import { ContractStatus } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { contractButtonSelect } from '@organization/src/components/Tabs/shared/contractSelectStyles';
import { ContractEndModal } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/EndContractModal';

import { contractOptionIcon } from './utils';

interface ContractStatusSelectProps {
  status: ContractStatus;
}

export const contractStatusOptions: SelectOption<ContractStatus>[] = [
  { label: 'Draft', value: ContractStatus.Draft },
  { label: 'Ended', value: ContractStatus.Ended },
  { label: 'Live', value: ContractStatus.Live },
];

export const ContractStatusSelect: React.FC<ContractStatusSelectProps> = ({
  status,
}) => {
  const { onOpen, onClose, isOpen } = useDisclosure({
    id: 'end-contract-modal',
  });
  const selected = contractStatusOptions.find((e) => e.value === status);
  const icon = contractOptionIcon?.[selected?.value];

  return (
    <>
      <Flex alignItems='center' gap={1} onClick={() => onOpen()}>
        {icon && (
          <Flex alignItems='center' boxSize={3}>
            {icon}
          </Flex>
        )}
        <Text
          color={status === ContractStatus.Live ? 'primary.800' : 'gray.800'}
        >
          {selected?.label}
        </Text>
      </Flex>
      <ContractEndModal
        isOpen={isOpen}
        onClose={onClose}
        contractId={'contractId'}
        organizationName={'organizationName'}
      />
    </>
  );
};
