import { FC, PropsWithChildren } from 'react';

import { DateTimeUtils } from '@utils/date';

interface TimelineItemProps extends PropsWithChildren {
  date: string;
  showDate: boolean;
}

export const TimelineItem: FC<TimelineItemProps> = ({
  date,
  showDate,
  children,
}) => {
  return (
    <div className='px-6 pb-2 bg-gray-25'>
      {showDate && (
        <span className='text-gray-500 text-xs font-medium mb-2 inline-block'>
          {DateTimeUtils.format(date, DateTimeUtils.defaultFormatShortString)}
        </span>
      )}
      {children}
    </div>
  );
};
