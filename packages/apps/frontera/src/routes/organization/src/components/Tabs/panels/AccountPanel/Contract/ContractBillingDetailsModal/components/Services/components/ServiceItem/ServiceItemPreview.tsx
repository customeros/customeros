import { FC } from 'react';

import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { DateTimeUtils } from '@utils/date.ts';
import { BilledType, ContractStatus } from '@graphql/types';
import { FlipBackward } from '@ui/media/icons/FlipBackward.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { currencySymbol } from '@shared/util/currencyOptions.ts';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  service: ContractLineItemStore;
  allowIndividualRestore?: boolean;
  type: 'subscription' | 'one-time';
  contractStatus?: ContractStatus | null;
}

const billedTypeLabel: Record<
  Exclude<BilledType, BilledType.None | BilledType.Usage | BilledType.Once>,
  string
> = {
  [BilledType.Monthly]: 'month',
  [BilledType.Quarterly]: 'quarter',
  [BilledType.Annually]: 'year',
};

const deleteButtonClasses =
  'border-none bg-transparent shadow-none text-gray-400 pr-3 pl-4 py-2 -mx-4 absolute -right-7 top-0 bottom-0 invisible group-hover:visible hover:bg-transparent';

export const ServiceItemPreview: FC<ServiceItemProps> = observer(
  ({
    service,
    isEnded,
    currency,
    type,
    contractStatus,
    allowIndividualRestore,
  }) => {
    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';

    const isFutureVersion =
      service?.tempValue?.serviceStarted &&
      DateTimeUtils.isFuture(service?.tempValue?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const isCurrentVersion =
      (service?.tempValue?.serviceEnded &&
        DateTimeUtils.isFuture(service?.tempValue?.serviceEnded) &&
        DateTimeUtils.isPast(service?.tempValue?.serviceStarted)) ||
      (!service?.tempValue?.serviceEnded &&
        DateTimeUtils.isPast(service?.tempValue?.serviceStarted));

    return (
      <>
        <div
          className={cn(
            'flex items-baseline justify-between group text-gray-700 relative',
            {
              'text-gray-400': isEnded,
              'line-through text-gray-400 hover:text-gray-400':
                service.tempValue.closed,
            },
          )}
        >
          <div className='flex items-baseline text-inherit'>
            <span>
              {service?.tempValue?.quantity}
              <span className='relative z-[2] mx-1'>×</span>

              {sliCurrencySymbol}
              {service?.tempValue?.price}
            </span>
            {type !== 'one-time' && <span>/ </span>}
            <span>
              {' '}
              {
                billedTypeLabel[
                  service?.tempValue?.billingCycle as Exclude<
                    BilledType,
                    BilledType.None | BilledType.Usage | BilledType.Once
                  >
                ]
              }{' '}
            </span>

            <span className='ml-1 text-inherit'>
              • {service?.tempValue?.tax?.taxRate}% VAT
            </span>
          </div>

          <div className='ml-1 text-inherit'>
            {isCurrentVersion && 'Current since '}

            {service?.tempValue?.serviceStarted &&
              DateTimeUtils.format(
                toZonedTime(
                  service?.tempValue?.serviceStarted,
                  'UTC',
                ).toString(),
                DateTimeUtils.dateWithShortYear,
              )}
          </div>
          {allowIndividualRestore &&
            (!service?.tempValue?.metadata.id || isDraft || isFutureVersion) &&
            service?.tempValue?.closed && (
              <IconButton
                size='xs'
                variant='outline'
                aria-label={'Restore version'}
                className={deleteButtonClasses}
                icon={<FlipBackward className='text-inherit' />}
                onClick={() =>
                  service.update((prev) => ({ ...prev, closed: false }), {
                    mutate: false,
                  })
                }
              />
            )}
        </div>
      </>
    );
  },
);
