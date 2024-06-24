import { DateTimeUtils } from '@utils/date.ts';
import { ContractStatus } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { contractOptionIcon } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractCardActions/utils.tsx';

export const ContractStatusTag = ({
  status,
  contractStarted,
}: {
  status: ContractStatus;
  contractStarted?: string;
}) => {
  const statusColorScheme: Record<string, string> = {
    [ContractStatus.Live]: 'primary',
    [ContractStatus.Draft]: 'gray',
    [ContractStatus.Ended]: 'gray',
    [ContractStatus.Scheduled]: 'primary',
    [ContractStatus.OutOfContract]: 'warning',
  };
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
  const icon = contractOptionIcon?.[status];
  const selected = contractStatusOptions.find((e) => e.value === status);

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
