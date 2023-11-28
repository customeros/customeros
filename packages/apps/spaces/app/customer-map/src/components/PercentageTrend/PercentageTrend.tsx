'use client';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { TrendUp01 } from '@ui/media/icons/TrendUp01';
import { TrendDown01 } from '@ui/media/icons/TrendDown01';

export const PercentageTrend = ({ percentage }: { percentage: number }) => (
  <Flex align='center' gap='1'>
    {percentage > 0 ? (
      <TrendUp01 boxSize='5' color='green.500' />
    ) : (
      <TrendDown01 boxSize='5' color='yellow.500' />
    )}
    <Text fontSize='sm' color={percentage > 0 ? 'green.600' : 'yellow.600'}>
      {Math.abs(percentage)}%
    </Text>
    <Text fontSize='sm'>vs last mth</Text>
  </Flex>
);
