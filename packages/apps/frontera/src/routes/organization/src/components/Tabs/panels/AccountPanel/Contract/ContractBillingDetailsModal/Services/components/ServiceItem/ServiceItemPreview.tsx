import { FC } from 'react';

import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { DateTimeUtils } from '@spaces/utils/date.ts';
import { BilledType, ContractStatus } from '@graphql/types';
import { FlipBackward } from '@ui/media/icons/FlipBackward.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { currencySymbol } from '@shared/util/currencyOptions.ts';

import ServiceLineItemStore from '../../../stores/Service.store.ts';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  service: ServiceLineItemStore;
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
      service?.serviceLineItem?.serviceStarted &&
      DateTimeUtils.isFuture(service?.serviceLineItem?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const isCurrentVersion =
      (service?.serviceLineItem?.serviceEnded &&
        DateTimeUtils.isFuture(service?.serviceLineItem?.serviceEnded) &&
        DateTimeUtils.isPast(service?.serviceLineItem?.serviceStarted)) ||
      (!service?.serviceLineItem?.serviceEnded &&
        DateTimeUtils.isPast(service?.serviceLineItem?.serviceStarted));

    return (
      <>
        <div
          className={cn(
            'flex items-baseline justify-between group text-gray-700 relative',
            {
              'text-gray-400': isEnded,
              'line-through text-gray-400 hover:text-gray-400':
                service.serviceLineItem?.isDeleted,
            },
          )}
        >
          <div className='flex items-baseline text-inherit'>
            <span>
              {service?.serviceLineItem?.quantity}
              <span className='relative z-[2] mx-1'>×</span>

              {sliCurrencySymbol}
              {service?.serviceLineItem?.price}
            </span>
            {type !== 'one-time' && <span>/ </span>}
            <span>
              {' '}
              {
                billedTypeLabel[
                  service?.serviceLineItem?.billingCycle as Exclude<
                    BilledType,
                    BilledType.None | BilledType.Usage | BilledType.Once
                  >
                ]
              }{' '}
            </span>

            <span className='ml-1 text-inherit'>
              • {service?.serviceLineItem?.tax?.taxRate}% VAT
            </span>
          </div>

          <div className='ml-1 text-inherit'>
            {isCurrentVersion && 'Current since '}

            {service?.serviceLineItem?.serviceStarted &&
              DateTimeUtils.format(
                toZonedTime(
                  service.serviceLineItem.serviceStarted,
                  'UTC',
                ).toString(),
                DateTimeUtils.dateWithShortYear,
              )}
          </div>
          {allowIndividualRestore &&
            (service.serviceLineItem?.isNew || isDraft || isFutureVersion) &&
            service.serviceLineItem?.isDeleted && (
              <IconButton
                aria-label={'Restore version'}
                icon={<FlipBackward className='text-inherit' />}
                variant='outline'
                size='xs'
                onClick={() => service.setIsDeleted(false)}
                className={deleteButtonClasses}
              />
            )}
        </div>
      </>
    );
  },
);
