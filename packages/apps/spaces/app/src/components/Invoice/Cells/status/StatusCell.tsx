import { Flex } from '@ui/layout/Flex';
import { Clock } from '@ui/media/icons/Clock';
import { InvoiceStatus } from '@graphql/types';
import { Tag, TagLeftIcon } from '@ui/presentation/Tag';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';

interface StatusCellProps {
  status?: InvoiceStatus | null;
}
function renderStatusNode(type: InvoiceStatus | null | undefined) {
  switch (type) {
    // case 'SCHEDULED':
    //   return (
    //     <Tag colorScheme='gray' variant='outline'>
    //       <TagLeftIcon as={ClockFastForward} />
    //       Scheduled
    //     </Tag>
    //   );
    case InvoiceStatus.Draft:
      return (
        <Tag colorScheme='gray' variant='outline'>
          <TagLeftIcon as={ClockFastForward} />
          Draft
        </Tag>
      );
    case InvoiceStatus.Paid:
      return (
        <Tag colorScheme='success' variant='outline'>
          <TagLeftIcon as={CheckCircle} />
          Paid
        </Tag>
      );
    // case 'PARTIALLY_PAID':
    //   return (
    //     <Tag colorScheme='success' variant='outline'>
    //       <TagLeftIcon as={CheckCircle} />
    //       Partially paid
    //     </Tag>
    //   );
    // case 'OVERDUE':
    //   return (
    //     <Tag colorScheme='warning' variant='outline'>
    //       <TagLeftIcon as={AlertCircle} />
    //       Overdue
    //     </Tag>
    //   );
    case 'DUE':
      return (
        <Tag colorScheme='primary' variant='outline'>
          <TagLeftIcon as={Clock} />
          Due
        </Tag>
      );
    case 'VOID':
      return (
        <Tag colorScheme='gray' variant='outline'>
          <TagLeftIcon as={SlashCircle01} />
          Voided
        </Tag>
      );
    default:
      return null;
  }
}

export const StatusCell = ({ status }: StatusCellProps) => {
  return <Flex align='center'>{renderStatusNode(status)}</Flex>;
};
