import { Clock } from '@ui/media/icons/Clock';
import { InvoiceStatus } from '@graphql/types';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag/Tag';

interface StatusCellProps {
  className?: string;
  status?: InvoiceStatus | null;
}

export function renderStatusNode(type: InvoiceStatus | null | undefined) {
  switch (type) {
    case InvoiceStatus.Initialized:
      return (
        <Tag variant='outline' colorScheme='gray'>
          <TagLeftIcon>
            <ClockFastForward />
          </TagLeftIcon>
          <TagLabel>Draft</TagLabel>
        </Tag>
      );
    case InvoiceStatus.Paid:
      return (
        <Tag variant='outline' colorScheme='success'>
          <TagLeftIcon>
            <CheckCircle />
          </TagLeftIcon>
          <TagLabel>Paid</TagLabel>
        </Tag>
      );
    case InvoiceStatus.Due:
      return (
        <Tag variant='outline' colorScheme='primary'>
          <TagLeftIcon>
            <Clock />
          </TagLeftIcon>
          <TagLabel>Due</TagLabel>
        </Tag>
      );
    case InvoiceStatus.Void:
      return (
        <Tag variant='outline' colorScheme='gray'>
          <TagLeftIcon>
            <SlashCircle01 />
          </TagLeftIcon>
          <TagLabel>Voided</TagLabel>
        </Tag>
      );
    case InvoiceStatus.Scheduled:
      return (
        <Tag variant='outline' colorScheme='gray'>
          <TagLeftIcon>
            <ClockFastForward />
          </TagLeftIcon>
          <TagLabel>Scheduled</TagLabel>
        </Tag>
      );
    case InvoiceStatus.Overdue:
      return (
        <Tag variant='outline' colorScheme='warning'>
          <TagLeftIcon>
            <InfoCircle />
          </TagLeftIcon>
          <TagLabel>Overdue</TagLabel>
        </Tag>
      );
    default:
      return null;
  }
}

export const StatusCell = ({ status }: StatusCellProps) => {
  return <div className='flex items-center'>{renderStatusNode(status)}</div>;
};
