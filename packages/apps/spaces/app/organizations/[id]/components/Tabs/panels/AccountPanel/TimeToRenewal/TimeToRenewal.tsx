import React from 'react';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { BillingDetails as BillingDetailsT } from '@graphql/types';
import { Text } from '@chakra-ui/react';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { getTimeToNextRenewal } from '../../../shared/util';

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
  const [numberValue, unit] = getTimeToNextRenewal(
    new Date(renewalCycleStart),
    renewalCycle,
  );
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
