import React from 'react';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { BillingDetails as BillingDetailsT } from '@graphql/types';
import { Text } from '@chakra-ui/react';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import {
  addWeeks,
  addMonths,
  addYears,
  differenceInWeeks,
  differenceInMonths,
  differenceInYears,
  formatDistanceToNowStrict,
} from 'date-fns';

type RenewalFrequency =
  | 'WEEKLY'
  | 'BIWEEKLY'
  | 'MONTHLY'
  | 'QUARTERLY'
  | 'BIANNUALLY'
  | 'ANNUALLY';
function getTimeToNextRenewal(
  renewalStart: Date,
  renewalFrequency: RenewalFrequency,
): string[] {
  const now = new Date();
  let totalRenewals: number;
  let nextRenewalDate: Date;

  switch (renewalFrequency) {
    case 'WEEKLY':
      totalRenewals = differenceInWeeks(now, renewalStart);
      nextRenewalDate = addWeeks(renewalStart, totalRenewals + 1);
      break;
    case 'BIWEEKLY':
      totalRenewals = differenceInWeeks(now, renewalStart) / 2;
      nextRenewalDate = addWeeks(renewalStart, 2 * (totalRenewals + 1));
      break;
    case 'MONTHLY':
      totalRenewals = differenceInMonths(now, renewalStart);
      nextRenewalDate = addMonths(renewalStart, totalRenewals + 1);
      break;
    case 'QUARTERLY':
      totalRenewals = differenceInMonths(now, renewalStart) / 3;
      nextRenewalDate = addMonths(renewalStart, 3 * (totalRenewals + 1));
      break;
    case 'BIANNUALLY':
      totalRenewals = differenceInMonths(now, renewalStart) / 6;
      nextRenewalDate = addMonths(renewalStart, 6 * (totalRenewals + 1));
      break;
    case 'ANNUALLY':
      totalRenewals = differenceInYears(now, renewalStart);
      nextRenewalDate = addYears(renewalStart, totalRenewals + 1);
      break;
    default:
      throw new Error('Unrecognized renewal frequency');
  }

  const distanceToNextRenewal = formatDistanceToNowStrict(nextRenewalDate, {
  });

  return distanceToNextRenewal?.split(' ');
}

export type BillingDetailsType = BillingDetailsT & { amount?: string | null };
interface BillingDetailsCardBProps {
  renewalCycleStart: BillingDetailsType['renewalCycleStart'];
  renewalCycle: BillingDetailsType['renewalCycle'];
}
export const TimeToRenewal: React.FC<BillingDetailsCardBProps> = ({
  renewalCycle,
  renewalCycleStart,
}) => {
  if (!renewalCycle || !renewalCycleStart) return null;
  const [numberValue, unit] = getTimeToNextRenewal(new Date(renewalCycleStart), renewalCycle);
  return (
    <Flex
      width='full'
      p={4}
      borderRadius='xl'
      border='1px solid'
      borderColor='gray.200'
      boxShadow='xs'
      justifyContent='space-between'
      bg='white'
    >
      <Flex alignItems='center'>
        <FeaturedIcon>
          <Icons.ClockFastForward />
        </FeaturedIcon>
        <Heading size='sm' color='gray.700' ml={4}>
          Time to renewal
        </Heading>
      </Flex>

      <Flex direction='column' alignItems='flex-end' justifyItems='center'>
        <Text fontSize='2xl' fontWeight='bold' lineHeight='1' color='gray.700'>
          {numberValue}
        </Text>
        <Text color='gray.500'>{unit}</Text>
      </Flex>
    </Flex>
  );
};
