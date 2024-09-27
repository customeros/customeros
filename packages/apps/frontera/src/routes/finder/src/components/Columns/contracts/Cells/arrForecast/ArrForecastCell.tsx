import { useMemo, useState } from 'react';

import { twMerge } from 'tailwind-merge';
import { observer } from 'mobx-react-lite';
import { PopoverTrigger } from '@radix-ui/react-popover';

import { cn } from '@ui/utils/cn';
import { ContractStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
} from '@ui/overlay/Popover/Popover.tsx';
import {
  RangeSlider,
  RangeSliderThumb,
  RangeSliderTrack,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider';

interface OwnerProps {
  id: string;
}

export const ArrForecastCell = observer(({ id }: OwnerProps) => {
  const store = useStore();
  const contract = store.contracts.value.get(id);
  const [isEditing, setIsEditing] = useState(false);
  const currency = 'USD';
  const opportunityId = useMemo(() => {
    return (
      contract?.value?.opportunities?.find((e) => e.internalStage === 'OPEN')
        ?.id || contract?.value?.opportunities?.[0]?.id
    );
  }, []);

  const opportunity = store.opportunities.value.get(opportunityId ?? '');
  const value = opportunity?.value?.renewalAdjustedRate;

  const trackStyle = cn('h-0.5 transition-colors', {
    'bg-orangeDark-700': value <= 25,
    'bg-yellow-400': value > 25 && value < 75,
    'bg-greenLight-400': value >= 75,
  });

  const thumbStyle = cn('ring-1 transition-colors shadow-md cursor-pointer', {
    'ring-orangeDark-700': value <= 25,
    'ring-yellow-400': value > 25 && value < 75,
    'ring-greenLight-400': value >= 75,
  });

  const hasEnded = contract?.value?.contractStatus === ContractStatus.Ended;

  const formattedAmount = formatCurrency(
    hasEnded ? 0 : opportunity?.value?.amount ?? 0,
    2,
    currency || 'USD',
  );
  const showPotentialAmount =
    opportunity?.value?.amount !== null &&
    opportunity?.value?.maxAmount !== null &&
    (opportunity?.value?.maxAmount ?? 0) * (value / 100) !==
      opportunity?.value?.maxAmount;

  if (!opportunity?.value?.amount && !opportunity?.value?.maxAmount)
    return <span className='text-gray-400'>Unknown</span>;

  const textColor = opportunity?.value?.amount
    ? 'text-gray-700'
    : 'text-gray-500';

  return (
    <div className='flex flex-col justify-center group/forecast'>
      <Popover open={isEditing} onOpenChange={setIsEditing}>
        <div className='flex gap-1 items-center'>
          <PopoverAnchor>
            <span className={twMerge('text-sm', textColor)}>
              {formattedAmount}
            </span>
          </PopoverAnchor>

          <PopoverTrigger asChild>
            {/*<IconButton*/}
            {/*  size='xxs'*/}
            {/*  variant='ghost'*/}
            {/*  aria-label='edit renewal likelihood'*/}
            {/*  icon={<Edit03 className='text-gray-500' />}*/}
            {/*  className={cn(*/}
            {/*    'rounded-md opacity-0 group-hover/forecast:opacity-100',*/}
            {/*    isEditing && 'opacity-100',*/}
            {/*  )}*/}
            {/*/>*/}
          </PopoverTrigger>
        </div>

        <PopoverContent
          align='start'
          sideOffset={showPotentialAmount ? 30 : 20}
        >
          <RangeSlider
            min={0}
            step={1}
            max={100}
            value={[value]}
            // className='w-40'
            // onValueChange={(values) => {
            //   setValue(values[0]);
            // }}
            // onValueCommit={(values) => {
            //   handleChange(values[0]);
            // }}
            // onKeyDown={(e) => {
            //   if (e.key === 'Escape') {
            //     setValue(adjustedRate ?? 0);
            //   }
            // }}
          >
            <RangeSliderTrack className='bg-gray-400 h-0.5'>
              <RangeSliderFilledTrack className={trackStyle} />
            </RangeSliderTrack>
            <RangeSliderThumb className={thumbStyle} />
          </RangeSlider>
        </PopoverContent>
      </Popover>
    </div>
  );
});
