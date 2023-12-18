'use client';
import sample from 'lodash/sample';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Minus } from '@ui/media/icons/Minus';
import { TrendUp01 } from '@ui/media/icons/TrendUp01';
import { TrendDown01 } from '@ui/media/icons/TrendDown01';

const quotes = [
  'As stable as a turtle on tranquilizers',
  'Marching to the rhythm of a metronome',
  'Gliding on the buttery highway of sameness',
  'Like a hamster on a wheel of monotony',
  "As consistent as a penguin's waddle",
  'Stuck in the syrup of predictability',
  "As routine as a robot's tea party",
  'Riding the carousel of constancy',
  'As predictable as a weather forecast in Arizona',
  "Like a koala's enthusiasm for change",
  'Chasing the tail of uniformity',
  'Dancing to the tune of déjà vu',
];

export const PercentageTrend = ({
  percentage,
}: {
  percentage: string | number;
}) => {
  percentage = `${percentage}`;
  const icon =
    percentage.indexOf('0') == 0 ? (
      <Minus boxSize='5' color='gray.700' />
    ) : percentage.indexOf('-') == 0 ? (
      <TrendDown01 boxSize='5' color='yellow.500' />
    ) : (
      <TrendUp01 boxSize='5' color='green.500' />
    );

  const color =
    percentage.indexOf('0') == 0
      ? 'gray.700'
      : percentage.indexOf('-') == 0
      ? 'yellow.600'
      : 'green.600';

  const quote = percentage.indexOf('0') == 0 ? sample(quotes) : 'vs last mth';

  return (
    <Flex align='center' gap='1'>
      {icon}
      <Text fontSize='sm' color={color}>
        {percentage}
      </Text>
      <Text fontSize='sm'>{quote}</Text>
    </Flex>
  );
};
