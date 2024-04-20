import React, { ReactNode } from 'react';

import { ContractStatus } from '@graphql/types';
import { DateTimeUtils } from '@spaces/utils/date';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag/Tag';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { contractOptionIcon } from './utils';

interface ContractStatusSelectProps {
  status: ContractStatus;
  statusContent: ReactNode;
  contractStarted?: string;
  onHandleStatusChange: () => void;
}

const statusColorScheme: Record<string, string> = {
  [ContractStatus.Live]: 'primary',
  [ContractStatus.Draft]: 'gray',
  [ContractStatus.Ended]: 'gray',
  [ContractStatus.Scheduled]: 'primary',
  [ContractStatus.OutOfContract]: 'warning',
};

export const ContractStatusTag: React.FC<ContractStatusSelectProps> = ({
  status,
  contractStarted,
  statusContent,
  onHandleStatusChange,
}) => {
  const contractStatusOptions: SelectOption<ContractStatus>[] = [
    { label: 'Draft', value: ContractStatus.Draft },
    { label: 'Ended', value: ContractStatus.Ended },
    { label: 'Live', value: ContractStatus.Live },
    { label: 'Out of contract', value: ContractStatus.OutOfContract },
    {
      label: contractStarted
        ? `Live ${DateTimeUtils.format(
            contractStarted,
            DateTimeUtils.defaultFormatShortString,
          )}`
        : 'Scheduled',
      value: ContractStatus.Scheduled,
    },
  ];

  const selected = contractStatusOptions.find((e) => e.value === status);
  const icon = contractOptionIcon?.[status];

  return (
    <>
      <Menu>
        <MenuButton disabled={status === ContractStatus.Scheduled}>
          <Tag
            className='flex items-center gap-1 whitespace-nowrap mx-0 px-1'
            colorScheme={statusColorScheme[status] as 'primary'}
          >
            <TagLeftIcon className='m-0'>{icon}</TagLeftIcon>

            <TagLabel>{selected?.label}</TagLabel>
          </Tag>
        </MenuButton>

        <MenuList align='end' side='bottom'>
          <MenuItem
            className='flex items-center'
            onClick={onHandleStatusChange}
          >
            {statusContent}
          </MenuItem>
        </MenuList>
      </Menu>
    </>
  );
};
