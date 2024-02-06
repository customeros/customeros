import { Text } from '@ui/typography/Text';

import { getMilestoneDueDate } from './utils';

interface PlanDueDateProps {
  value: string;
  isDone: boolean;
}

export const PlanDueDate = ({ value, isDone }: PlanDueDateProps) => {
  if (isDone) return null;

  return (
    <Text as='label' fontSize='sm' color='gray.500' whiteSpace='nowrap'>
      {getMilestoneDueDate(value, isDone).toLowerCase()}
    </Text>
  );
};
