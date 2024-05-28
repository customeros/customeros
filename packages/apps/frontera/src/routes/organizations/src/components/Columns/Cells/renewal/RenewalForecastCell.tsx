import { useState, useEffect } from 'react';

import set from 'lodash/set';
import { twMerge } from 'tailwind-merge';
import { observer } from 'mobx-react-lite';
import { PopoverTrigger } from '@radix-ui/react-popover';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';
import {
  RangeSlider,
  RangeSliderThumb,
  RangeSliderTrack,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider/RangeSlider';

interface RenewalForecastCellProps {
  id: string;
  amount?: number | null;
  potentialAmount?: number | null;
}

export const RenewalForecastCell = observer(
  ({ id }: RenewalForecastCellProps) => {
    const store = useStore();
    const [isEditing, setIsEditing] = useState(false);

    const organization = store.organizations.value.get(id);
    const amount =
      organization?.value.accountDetails?.renewalSummary?.arrForecast ?? null;
    const potentialAmount =
      organization?.value.accountDetails?.renewalSummary?.maxArrForecast ??
      null;

    const initialValue = (() => {
      if (potentialAmount === 0) return 0;

      return ((amount ?? 0) / (potentialAmount ?? 0)) * 100;
    })();

    const [value, setValue] = useState(initialValue);

    const formattedAmount =
      amount !== null && amount >= 0
        ? formatCurrency((potentialAmount ?? 0) * (value / 100), 0)
        : 'Unknown';
    const formattedPotentialAmount = formatCurrency(potentialAmount ?? 0, 0);

    const showPotentialAmount =
      amount !== null &&
      potentialAmount !== null &&
      (potentialAmount ?? 0) * (value / 100) !== potentialAmount;

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

    const handleChange = (value: number) => {
      store.organizations.value.get(id)?.update((org) => {
        set(
          org,
          'accountDetails.renewalSummary.renewalLikelihood',
          (() => {
            if (value <= 25) return OpportunityRenewalLikelihood.LowRenewal;
            if (value > 25 && value < 75)
              return OpportunityRenewalLikelihood.MediumRenewal;

            return OpportunityRenewalLikelihood.HighRenewal;
          })(),
        );
        set(
          org,
          'accountDetails.renewalSummary.arrForecast',
          (potentialAmount ?? 0) * (value / 100),
        );

        return org;
      });
    };

    useEffect(() => {
      setValue(initialValue);
    }, [initialValue]);

    if (formattedAmount === 'Unknown')
      return <span className='text-gray-400'>Unknown</span>;
    const textColor = amount ? 'text-gray-700' : 'text-gray-500';

    return (
      <div className='flexjustify-start group/forecast'>
        <Popover open={isEditing} onOpenChange={setIsEditing}>
          <div className='flex gap-1 items-center'>
            <PopoverAnchor>
              <span>
                {showPotentialAmount && (
                  <span className='text-sm text-gray-500 line-through mr-1'>
                    {formattedPotentialAmount}
                  </span>
                )}
                <span className={twMerge('text-sm', textColor)}>
                  {formattedAmount}
                </span>
              </span>
            </PopoverAnchor>

            <PopoverTrigger asChild>
              <IconButton
                size='xxs'
                variant='ghost'
                aria-label='edit renewal likelihood'
                icon={<Edit03 className='text-gray-500' />}
                className={cn(
                  'rounded-md opacity-0 group-hover/forecast:opacity-100',
                  isEditing && 'opacity-100',
                )}
              />
            </PopoverTrigger>
          </div>

          <PopoverContent
            sideOffset={showPotentialAmount ? 30 : 20}
            align='start'
          >
            <RangeSlider
              step={1}
              min={0}
              max={100}
              value={[value]}
              className='w-40'
              onKeyDown={(e) => {
                if (e.key === 'Escape') {
                  setValue(initialValue);
                }
              }}
              onValueChange={(values) => {
                setValue(values[0]);
              }}
              onValueCommit={(values) => {
                handleChange(values[0]);
              }}
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
  },
);
