import { useState, useEffect } from 'react';

import set from 'lodash/set';
import { produce } from 'immer';
import { twMerge } from 'tailwind-merge';
import { useQueryClient } from '@tanstack/react-query';
import { PopoverTrigger } from '@radix-ui/react-popover';
import { useInfiniteGetRenewalsQuery } from '@renewals/graphql/getRenewals.generated';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { toastError } from '@ui/presentation/Toast';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { OpportunityRenewalLikelihood } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useRenewalsMeta } from '@shared/state/RenewalsMeta.atom';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';
import { useUpdateOpportunityRenewalMutation } from '@organization/graphql/updateOpportunityRenewal.generated';
import {
  RangeSlider,
  RangeSliderThumb,
  RangeSliderTrack,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider/RangeSlider';

interface RenewalForecastCellProps {
  id: string;
  opportunityId: string;
  amount?: number | null;
  adjustedRate?: number | null;
  potentialAmount?: number | null;
}

export const RenewalForecastCell = ({
  id,
  amount = null,
  opportunityId,
  adjustedRate = 0,
  potentialAmount = null,
}: RenewalForecastCellProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [renewalsMeta] = useRenewalsMeta();

  const [value, setValue] = useState(adjustedRate ?? 0);

  const { getRenewals } = renewalsMeta;
  const queryKey = useInfiniteGetRenewalsQuery.getKey(getRenewals);

  const updateOpportunityRenewal = useUpdateOpportunityRenewalMutation(client, {
    onMutate: () => {
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = useInfiniteGetRenewalsQuery.mutateCacheEntry(
        queryClient,
        getRenewals,
      )((old) =>
        produce(old, (draft) => {
          const pageIndex = renewalsMeta.getRenewals.pagination.page - 1;

          const content =
            draft?.pages?.[pageIndex]?.dashboardView_Renewals?.content;
          const index = content?.findIndex(
            (item) => item.organization.metadata.id === id,
          );

          const nextLikelihood = (() => {
            if (value <= 25) return OpportunityRenewalLikelihood.LowRenewal;
            if (value > 25 && value < 75)
              return OpportunityRenewalLikelihood.MediumRenewal;

            return OpportunityRenewalLikelihood.HighRenewal;
          })();

          if (content && index !== undefined && index > -1) {
            set(
              content[index],
              'opportunity.renewalLikelihood',
              nextLikelihood,
            );
            set(
              content[index],
              'opportunity.amount',
              (potentialAmount ?? 0) * (value / 100),
            );
          }
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      toastError(
        `We couldn't update the forecast`,
        'renewal-forecast-update-error',
      );

      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
    },
    onSettled: () => {
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
      }, 500);
    },
  });

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
    updateOpportunityRenewal.mutate({
      input: {
        opportunityId,
        renewalAdjustedRate: value,
        renewalLikelihood: (() => {
          if (value <= 25) return OpportunityRenewalLikelihood.LowRenewal;
          if (value > 25 && value < 75)
            return OpportunityRenewalLikelihood.MediumRenewal;

          return OpportunityRenewalLikelihood.HighRenewal;
        })(),
      },
    });
  };

  useEffect(() => {
    setValue(adjustedRate ?? 0);
  }, [adjustedRate]);

  if (!amount && !potentialAmount)
    return <span className='text-gray-400'>Unknown</span>;

  const textColor = amount ? 'text-gray-700' : 'text-gray-500';

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
          align='start'
          sideOffset={showPotentialAmount ? 30 : 20}
        >
          <RangeSlider
            min={0}
            step={1}
            max={100}
            value={[value]}
            className='w-40'
            onValueChange={(values) => {
              setValue(values[0]);
            }}
            onValueCommit={(values) => {
              handleChange(values[0]);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                setValue(adjustedRate ?? 0);
              }
            }}
          >
            <RangeSliderTrack className='bg-gray-400 h-0.5'>
              <RangeSliderFilledTrack className={trackStyle} />
            </RangeSliderTrack>
            <RangeSliderThumb className={thumbStyle} />
          </RangeSlider>
        </PopoverContent>
      </Popover>
      {showPotentialAmount && (
        <span className='text-sm text-gray-500 line-through'>
          {formattedPotentialAmount}
        </span>
      )}
    </div>
  );
};
