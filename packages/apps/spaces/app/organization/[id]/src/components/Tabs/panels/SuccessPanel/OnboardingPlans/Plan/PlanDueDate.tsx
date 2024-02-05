import differenceInDays from 'date-fns/differenceInDays';

import { Text } from '@ui/typography/Text';
interface PlanDueDateProps {
  value: string;
  isDone: boolean;
}

export const PlanDueDate = ({ value, isDone }: PlanDueDateProps) => {
  if (isDone) return null;

  const days = differenceInDays(new Date(value), new Date());

  const displayText = (() => {
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
