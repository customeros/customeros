import { Text } from '@ui/typography/Text';
import { getDifferenceFromNow } from '@shared/util/date';

interface TimeToRenewalCellProps {
  nextRenewalDate: string;
}

export const TimeToRenewalCell = ({
  nextRenewalDate,
}: TimeToRenewalCellProps) => {
  if (!nextRenewalDate)
    return (
      <Text fontSize='sm' color='gray.400'>
        Unknown
      </Text>
    );

  const [value, unit] = getDifferenceFromNow(nextRenewalDate);

  return (
    <Text fontSize='sm' color='gray.700'>
      {value} {unit}
    </Text>
  );
};
