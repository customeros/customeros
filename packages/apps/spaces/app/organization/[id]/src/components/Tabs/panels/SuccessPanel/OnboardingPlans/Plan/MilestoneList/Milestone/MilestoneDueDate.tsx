import differenceInDays from 'date-fns/differenceInDays';

import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';

interface MilestoneDueDateProps {
  value: string;
  isDone?: boolean;
}

export const MilestoneDueDate = ({ value, isDone }: MilestoneDueDateProps) => {
  const days = differenceInDays(new Date(value), new Date());
  const displayText = (() => {
    if (isDone) {
      return days !== 0
        ? `Done on ${DateTimeUtils.format(value, DateTimeUtils.date)}`
        : 'Done today';
    }

    const suffix = days === 1 ? 'day' : 'days';
    const prefix = days < 0 ? 'Late by' : days === 0 ? 'Due' : 'Due in';
    const absDays = Math.abs(days);

    if (absDays === 0) return 'Due today';

    return `${prefix} ${absDays} ${suffix}`;
  })();

  return (
    <Text as='label' fontSize='sm' color='gray.500' whiteSpace='nowrap'>
      {displayText}
    </Text>
  );
};
