import differenceInDays from 'date-fns/differenceInDays';

import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';

interface PlanDueDateProps {
  value: string;
  isDone: boolean;
}

export const PlanDueDate = ({ value, isDone }: PlanDueDateProps) => {
  const days = differenceInDays(new Date(value), new Date());

  const displayText = (() => {
    if (isDone) {
      return `Done on ${DateTimeUtils.format(value, DateTimeUtils.date)}`;
    }

    const suffix = days === 1 ? 'day' : 'days';
    const prefix = days < 0 ? 'late by' : days === 0 ? 'due' : 'due in';

    const absDays = Math.abs(days);

    if (absDays === 0) return 'due today';

    return `${prefix} ${absDays} ${suffix}`;
  })();

  return (
    <Text as='label' fontSize='sm' color='gray.500' whiteSpace='nowrap'>
      {displayText}
    </Text>
  );
};
