import { Text } from '@ui/typography/Text';
import { RenewalCycle } from '@graphql/types';

import { getTimeToRenewal } from '@organization/src/components/Tabs/shared/util';

interface TimeToRenewalCellProps {
  renewalDate: string | null;
  renewalFrequency?: RenewalCycle | null;
}

export const TimeToRenewalCell = ({
  renewalDate,
  renewalFrequency,
}: TimeToRenewalCellProps) => {
  if (!renewalDate || !renewalFrequency)
    return (
      <Text fontSize='sm' color='gray.500'>
        Unknown
      </Text>
    );

  const [value, unit] = getTimeToRenewal(renewalDate);

  return (
    <Text fontSize='sm' color='gray.700'>
      {value} {unit}
    </Text>
  );
};
