import { Text } from '@ui/typography/Text';
import { getTimeToRenewal } from '@organization/src/components/Tabs/shared/util';

interface TimeToRenewalCellProps {
  nextRenewalDate: string;
}

export const TimeToRenewalCell = ({
  nextRenewalDate,
}: TimeToRenewalCellProps) => {
  if (!nextRenewalDate)
    return (
      <Text fontSize='sm' color='gray.500'>
        Unknown
      </Text>
    );

  const [value, unit] = getTimeToRenewal(nextRenewalDate);

  return (
    <Text fontSize='sm' color='gray.700'>
      {value} {unit}
    </Text>
  );
};
