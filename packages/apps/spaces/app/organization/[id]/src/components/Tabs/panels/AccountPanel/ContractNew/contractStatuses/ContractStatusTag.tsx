import React from 'react';

import { ContractStatus } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag/Tag';

import { contractOptionIcon } from './utils';

interface ContractStatusSelectProps {
  status: ContractStatus;
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

export const ContractStatusTag: React.FC<ContractStatusSelectProps> = ({
  status,
}) => {
  const selected = contractStatusOptions.find((e) => e.value === status);
  const icon = contractOptionIcon?.[status];

  return (
    <>
      <Tag
        className='flex items-center gap-1 whitespace-nowrap mx-0 px-1'
        colorScheme={statusColorScheme[status] as 'primary'}
      >
        <TagLeftIcon className='m-0'>{icon}</TagLeftIcon>

        <TagLabel>{selected?.label}</TagLabel>
      </Tag>
    </>
  );
};
