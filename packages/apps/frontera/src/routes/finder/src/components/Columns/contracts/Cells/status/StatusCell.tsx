import { ContractStatus } from '@graphql/types';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface StatusCellProps {
  className?: string;
  status?: ContractStatus | null;
}

export function renderStatusNode(type: ContractStatus | null | undefined) {
  switch (type) {
    case ContractStatus.Draft:
      return (
        <Tag variant='subtle' colorScheme='grayModern'>
          <TagLabel>Draft</TagLabel>
        </Tag>
      );

    case ContractStatus.Live:
      return (
        <Tag variant='subtle' colorScheme='success'>
          <TagLabel>Live</TagLabel>
        </Tag>
      );

    case ContractStatus.Ended:
      return (
        <Tag variant='subtle' colorScheme='grayModern'>
          <TagLabel>Ended</TagLabel>
        </Tag>
      );

    case ContractStatus.OutOfContract:
      return (
        <Tag variant='subtle' colorScheme='warning'>
          <TagLabel>Out of contract</TagLabel>
        </Tag>
      );

    case ContractStatus.Scheduled:
      return (
        <Tag variant='subtle' colorScheme='primary'>
          <TagLabel>Scheduled</TagLabel>
        </Tag>
      );

    default:
      return '';
  }
}

export const StatusCell = ({ status }: StatusCellProps) => {
  return <div className='flex items-center'>{renderStatusNode(status)}</div>;
};
