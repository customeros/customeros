import { Text } from '@ui/typography/Text';

import { getMilestoneDueDate } from '../../utils';

interface MilestoneDueDateProps {
  value: string;
  isDone?: boolean;
}

export const MilestoneDueDate = ({ value, isDone }: MilestoneDueDateProps) => {
  return (
    <Text as='label' fontSize='sm' color='gray.500' whiteSpace='nowrap'>
      {getMilestoneDueDate(value, isDone)}
    </Text>
  );
};
