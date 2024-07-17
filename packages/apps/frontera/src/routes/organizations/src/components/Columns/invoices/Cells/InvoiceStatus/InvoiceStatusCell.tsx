import { InvoiceStatus } from '@graphql/types';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag/Tag';

export const InvoiceStatusCell = ({ status }: { status: InvoiceStatus }) => {
  return (
    <div className='flex flex-col items-start'>
      <Tag
        colorScheme={status === InvoiceStatus.OnHold ? 'warning' : 'grayBlue'}
      >
        <TagLeftIcon>
          <ClockFastForward />
        </TagLeftIcon>
        <TagLabel className='capitalize'>{status.toLowerCase()}</TagLabel>
      </Tag>
    </div>
  );
};
