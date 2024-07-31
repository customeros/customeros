import sample from 'lodash/sample';

import { cn } from '@ui/utils/cn';
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

export const PercentageTrend = ({ percentage }: { percentage: string }) => {
  percentage = `${percentage}`;

  const icon =
    percentage.indexOf('0') == 0 ? (
      <Minus className='size-5 text-gray-700' />
    ) : percentage.indexOf('-') == 0 ? (
      <TrendDown01 className='size-5 text-warning-500' />
    ) : (
      <TrendUp01 className='size-5 text-success-500' />
    );

  const color =
    percentage.indexOf('0') == 0
      ? 'text-gray-700'
      : percentage.indexOf('-') == 0
      ? 'text-yellow-600'
      : 'text-succes-600';

  const quote = percentage.indexOf('0') == 0 ? sample(quotes) : 'vs last mth';

  return (
    <div className='flex items-center gap-1'>
      {icon}
      <p className={cn(color, 'text-sm')}>{percentage}</p>
      <p className='text-sm'>{quote}</p>
    </div>
  );
};
