import { FC } from 'react';

import { Store } from '@store/store.ts';
import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { DateTimeUtils } from '@utils/date.ts';
import { FlipBackward } from '@ui/media/icons/FlipBackward.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { currencySymbol } from '@shared/util/currencyOptions.ts';
import { BilledType, ContractStatus, ServiceLineItem } from '@graphql/types';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  service: Store<ServiceLineItem>;
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
      service?.value?.serviceStarted &&
      DateTimeUtils.isFuture(service?.value?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const isCurrentVersion =
      (service?.value?.serviceEnded &&
        DateTimeUtils.isFuture(service?.value?.serviceEnded) &&
        DateTimeUtils.isPast(service?.value?.serviceStarted)) ||
      (!service?.value?.serviceEnded &&
        DateTimeUtils.isPast(service?.value?.serviceStarted));

    return (
      <>
        <div
          className={cn(
            'flex items-baseline justify-between group text-gray-700 relative',
            {
              'text-gray-400': isEnded,
              'line-through text-gray-400 hover:text-gray-400':
                service.value.closed,
            },
          )}
        >
          <div className='flex items-baseline text-inherit'>
            <span>
              {service?.value?.quantity}
              <span className='relative z-[2] mx-1'>×</span>

              {sliCurrencySymbol}
              {service?.value?.price}
            </span>
            {type !== 'one-time' && <span>/ </span>}
            <span>
              {' '}
              {
                billedTypeLabel[
                  service?.value?.billingCycle as Exclude<
                    BilledType,
                    BilledType.None | BilledType.Usage | BilledType.Once
                  >
                ]
              }{' '}
            </span>

            <span className='ml-1 text-inherit'>
              • {service?.value?.tax?.taxRate}% VAT
            </span>
          </div>

          <div className='ml-1 text-inherit'>
            {isCurrentVersion && 'Current since '}

            {service?.value?.serviceStarted &&
              DateTimeUtils.format(
                toZonedTime(service?.value?.serviceStarted, 'UTC').toString(),
                DateTimeUtils.dateWithShortYear,
              )}
          </div>
          {allowIndividualRestore &&
            (!service?.value?.metadata.id || isDraft || isFutureVersion) &&
            service?.value?.closed && (
              <IconButton
                aria-label={'Restore version'}
                icon={<FlipBackward className='text-inherit' />}
                variant='outline'
                size='xs'
                onClick={() =>
                  service.update((prev) => ({ ...prev, closed: false }), {
                    mutate: false,
                  })
                }
                className={deleteButtonClasses}
              />
            )}
        </div>
      </>
    );
  },
);
