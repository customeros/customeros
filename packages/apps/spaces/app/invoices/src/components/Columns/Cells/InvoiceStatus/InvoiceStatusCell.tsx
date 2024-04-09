import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag/Tag';

export const InvoiceStatusCell = ({
  isOutOfContract,
}: {
  isOutOfContract: boolean;
}) => {
  const label = isOutOfContract ? 'On hold' : 'Scheduled';

  return (
    <div className='flex flex-col items-start'>
      <Tag colorScheme={isOutOfContract ? 'warning' : 'grayBlue'}>
        <TagLeftIcon>
          <ClockFastForward />
        </TagLeftIcon>
        <TagLabel>{label}</TagLabel>
      </Tag>
    </div>
  );
};
