import { DateTimeUtils } from '@utils/date';

export const LastTouchpointDateCell = ({
  lastTouchPointAt,
}: {
  lastTouchPointAt: string;
}) => {
  return (
    <span className='text-gray-700'>
      {DateTimeUtils.timeAgo(lastTouchPointAt, {
        addSuffix: true,
      })}
    </span>
  );
};
