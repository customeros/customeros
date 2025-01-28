import { useState } from 'react';

import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store';

import { DateTimeUtils } from '@utils/date';
import { ContractStatus } from '@graphql/types';
import { Delete } from '@ui/media/icons/Delete';
import { toastError } from '@ui/presentation/Toast';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { PauseCircle } from '@ui/media/icons/PauseCircle.tsx';
import { currencySymbol } from '@shared/util/currencyOptions';
import { MaskedResizableInput } from '@ui/form/Input/MaskedResizableInput';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';

import { BilledTypeEditField } from './BilledTypeEditField';

interface ProductItemProps {
  currency?: string;
  isModification?: boolean;
  service: ContractLineItemStore;

  type: 'subscription' | 'one-time';
  allServices?: ContractLineItemStore[];

  contractStatus?: ContractStatus | null;
}

const inputClasses =
  'text-sm min-w-2.5 min-h-0 max-h-4 text-inherit underline hover:border-none focus:border-none border-none';

const deleteButtonClasses =
  'border-none bg-transparent shadow-none text-gray-400 pr-3 pl-4 py-2 -mx-4 absolute -right-7 top-0 bottom-0 invisible group-hover:visible hover:bg-transparent';

const formatPrice = (price: number | undefined | string): string => {
  if (price === undefined || price === null) return '';

  const priceStr = price.toString();

  // If no decimal point, return as is
  if (!priceStr.includes('.')) return `${priceStr}`;

  const [whole, decimal] = priceStr.split('.');

  if (!decimal) return `${whole}.`;

  // If decimal part exists and is single digit, add a zero
  if (decimal && decimal.length === 1) {
    return `${whole}.${decimal}0`;
  }

  return priceStr;
};

export const ProductItemEdit = observer(
  ({
    service,
    allServices,
    currency,
    isModification,
    type,
    contractStatus,
  }: ProductItemProps) => {
    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';
    const [displayPrice, setDisplayPrice] = useState(() =>
      formatPrice(service.tempValue.price),
    );
    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const onChangeServiceStarted = (e: Date | null) => {
      if (!e) return;

      const utcDate = toZonedTime(e, 'UTC');

      const checkExistingServiceStarted = (date: Date) => {
        return allServices?.some((val) => {
          const normalizedNew = `${date.toISOString().split('T')[0]}`;

          return (
            val.tempValue.metadata.id !== service.tempValue.metadata.id &&
            DateTimeUtils.isSameDay(val.tempValue.serviceStarted, normalizedNew)
          );
        });
      };

      const findCurrentService = () => {
        if (isDraft) return null;

        if (allServices?.length === 1) {
          return allServices[0]?.tempValue?.serviceStarted;
        }

        return allServices?.find((serviceData) => {
          const serviceStarted = serviceData?.tempValue?.serviceStarted;
          const serviceEnded = serviceData?.tempValue?.serviceEnded;

          return (
            (serviceEnded &&
              DateTimeUtils.isFuture(serviceEnded) &&
              DateTimeUtils.isPast(serviceStarted)) ||
            (!serviceEnded && DateTimeUtils.isPast(serviceStarted))
          );
        })?.tempValue?.serviceStarted;
      };

      const checkIfBeforeCurrentService = (
        date: Date,
        currentService: Date | null,
      ) => {
        if (allServices?.length === 1) {
          return false;
        }

        return (
          currentService &&
          DateTimeUtils.isBefore(date.toString(), currentService.toString())
        );
      };

      const checkIfBeforeToday = (date: Date) => {
        if (isDraft || type === 'one-time') return null;

        if (allServices?.length === 1) {
          return DateTimeUtils.isBefore(date.toString(), new Date().toString());
        }

        return false;
      };

      const existingServiceStarted = checkExistingServiceStarted(utcDate);
      const isTodayOrBefore = checkIfBeforeToday(utcDate);
      const currentService = findCurrentService();
      const isBeforeCurrentService = checkIfBeforeCurrentService(
        utcDate,
        currentService,
      );

      if (isTodayOrBefore) {
        toastError(
          `Select a service start date that is in the future`,
          `${service?.tempValue?.metadata?.id}-service-started-date-update-error`,
        );

        return;
      }

      if (isBeforeCurrentService) {
        toastError(
          `Modifications must be effective after the current service`,
          `${service?.tempValue?.metadata?.id}-service-started-date-update-error`,
        );

        return;
      }

      if (isBeforeCurrentService) {
        toastError(
          `Modifications must be effective after the current service`,
          `${service?.tempValue?.metadata?.id}-service-started-date-update-error`,
        );

        return;
      }

      if (existingServiceStarted) {
        toastError(
          `A version with this date already exists`,
          `${service?.tempValue?.metadata?.id}-service-started-date-update-error`,
        );

        return;
      }

      service.updateTemp((prev) => ({
        ...prev,
        serviceStarted: e.toISOString(),
      }));
    };

    const updateQuantity = (quantity: string) => {
      service.updateTemp((prev) => ({
        ...prev,
        quantity,
      }));
    };

    const updatePrice = (price: string) => {
      service.updateTemp(
        // @ts-expect-error  we allow undefined during edition but on blur we still enforce value therefore this is false positive
        (prev) => ({ ...prev, price: price ? price : undefined }),
      );
    };

    const updateTaxRate = (taxRate: string) => {
      service.updateTemp((prev) => ({
        ...prev,
        tax: {
          ...prev.tax,
          // @ts-expect-error we allow undefined during edition but on blur we still enforce value therefore this is false positive
          taxRate: taxRate ? parseFloat(taxRate) : undefined,
        },
      }));
    };

    return (
      <div className='flex items-baseline justify-between group relative text-gray-500 '>
        <div className='flex items-baseline'>
          <MaskedResizableInput
            min={0}
            size='xs'
            type='number'
            mask={Number}
            autofix={true}
            placeholder='0'
            className={inputClasses}
            onFocus={(e) => e.target.select()}
            value={service?.tempValue?.quantity?.toString() || ''}
            onChange={(e) => {
              updateQuantity(e.target.value ?? '');
            }}
            onBlur={(e) =>
              !e.target.value?.length
                ? updateQuantity('0')
                : updateQuantity(e.target.value)
            }
          />
          <span className=' mx-1 text-gray-700'>×</span>

          {sliCurrencySymbol}

          <MaskedResizableInput
            mask={`num`}
            placeholder='0'
            className={inputClasses}
            measureValue={displayPrice}
            onFocus={(e) => e.target.select()}
            value={service.tempValue.price?.toString() || ''}
            onAccept={(val, _d, e) => {
              updatePrice(val);

              if ((e?.target as HTMLInputElement)?.value) {
                setDisplayPrice(
                  formatPrice((e?.target as HTMLInputElement).value),
                );
              } else {
                setDisplayPrice(formatPrice(val));
              }
            }}
            blocks={{
              num: {
                mask: Number,
                scale: 2,
                radix: '.',
                lazy: true,
                min: type === 'one-time' ? -9999999999 : 0,
                placeholderChar: '#',
                thousandsSeparator: ',',
                normalizeZeros: true,
                padFractionalZeros: true,
                autofix: true,
              },
            }}
          />

          {type === 'one-time' ? (
            <span className='text-gray-700'></span>
          ) : (
            <BilledTypeEditField
              isModification={isModification}
              id={service.tempValue.metadata.id}
            />
          )}
          <span className=' mx-1 text-gray-700'>•</span>

          <MaskedResizableInput
            min={0}
            size='xs'
            type='number'
            mask={Number}
            autofix={true}
            placeholder='0'
            className={inputClasses}
            onFocus={(e) => e.target.select()}
            onChange={(e) => updateTaxRate(e.target.value)}
            value={service?.tempValue?.tax?.taxRate?.toString() || ''}
            onBlur={(e) =>
              !e.target.value?.trim()?.length
                ? updateTaxRate('0')
                : updateTaxRate(e.target.value)
            }
          />

          <span className='whitespace-nowrap  mx-1 text-gray-700'>% VAT</span>
        </div>
        <div className='flex items-center'>
          <Tooltip label='Service start date'>
            <div>
              <DatePickerUnderline
                onChange={onChangeServiceStarted}
                value={service?.tempValue?.serviceStarted}
              />
            </div>
          </Tooltip>
          <IconButton
            size='xs'
            variant='outline'
            aria-label={'Delete version'}
            className={deleteButtonClasses}
            icon={<Delete className='text-inherit' />}
            onClick={() => {
              service.updateTemp((prev) => ({ ...prev, closed: true }));
            }}
          />
          {service.tempValue.paused && (
            <Tooltip label={'This service will be invoiced when resumed'}>
              <PauseCircle className='text-gray-500 size-4 ml-2' />
            </Tooltip>
          )}
        </div>
      </div>
    );
  },
);
