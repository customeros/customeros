import { Text } from '@ui/typography/Text';
import { RenewalCycle } from '@spaces/graphql';

import {
  getTimeToNextRenewal,
  RenewalFrequency,
} from '@organization/components/Tabs/shared/util';

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

  const [numberValue, unit] = getTimeToNextRenewal(
    new Date(renewalDate ?? ''),
    renewalFrequency as RenewalFrequency,
  );

  return (
    <Text fontSize='sm' color='gray.700'>
      {numberValue} {unit}
    </Text>
  );
};
