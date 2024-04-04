import { Clock } from '@ui/media/icons/Clock';
import { InvoiceStatus } from '@graphql/types';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { SlashCircle01 } from '@ui/media/icons/SlashCircle01';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag/Tag';

interface StatusCellProps {
  status?: InvoiceStatus | null;
}
export function renderStatusNode(type: InvoiceStatus | null | undefined) {
  switch (type) {
    // case 'SCHEDULED':
    //   return (
    //     <Tag colorScheme='gray' variant='outline'>
    //       <TagLeftIcon>
    //         <ClockFastForward />
    //       </TagLeftIcon>
    //       <TagLabel>Scheduled</TagLabel>
    //     </Tag>
    //   );
    case InvoiceStatus.Initialized:
      return (
        <Tag colorScheme='gray' variant='outline'>
          <TagLeftIcon>
            <ClockFastForward />
          </TagLeftIcon>
          <TagLabel>Draft</TagLabel>
        </Tag>
      );
    case InvoiceStatus.Paid:
      return (
        <Tag colorScheme='success' variant='outline'>
          <TagLeftIcon>
            <CheckCircle />
          </TagLeftIcon>
          <TagLabel>Paid</TagLabel>
        </Tag>
      );
    // case 'PARTIALLY_PAID':
    //   return (
    //     <Tag colorScheme='success' variant='outline'>
    //       <TagLeftIcon>
    //         <CheckCircle />
    //       </TagLeftIcon>
    //       <TagLabel>Partially paid</TagLabel>
    //     </Tag>
    //   );
    // case 'OVERDUE':
    //   return (
    //     <Tag colorScheme='warning' variant='outline'>
    //       <TagLeftIcon>
    //         <AlertCircle />
    //       </TagLeftIcon>
    //       <TagLabel>Overdue</TagLabel>
    //     </Tag>
    //   );
    case 'DUE':
      return (
        <Tag colorScheme='primary' variant='outline'>
          <TagLeftIcon>
            <Clock />
          </TagLeftIcon>
          <TagLabel>Due</TagLabel>
        </Tag>
      );
    case 'VOID':
      return (
        <Tag colorScheme='gray' variant='outline'>
          <TagLeftIcon>
            <SlashCircle01 />
          </TagLeftIcon>
          <TagLabel>Voided</TagLabel>
        </Tag>
      );
    default:
      return null;
  }
}

export const StatusCell = ({ status }: StatusCellProps) => {
  return <div className='flex items-center'>{renderStatusNode(status)}</div>;
};
