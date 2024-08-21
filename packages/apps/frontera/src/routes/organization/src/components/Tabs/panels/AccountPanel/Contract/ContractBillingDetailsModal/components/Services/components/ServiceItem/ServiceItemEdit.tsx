import { observer } from 'mobx-react-lite';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { DateTimeUtils } from '@utils/date.ts';
import { ContractStatus } from '@graphql/types';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { toastError } from '@ui/presentation/Toast';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { currencySymbol } from '@shared/util/currencyOptions.ts';
import { MaskedResizableInput } from '@ui/form/Input/MaskedResizableInput.tsx';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2.tsx';

import { BilledTypeEditField } from './BilledTypeEditField.tsx';

interface ServiceItemProps {
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

export const ServiceItemEdit = observer(
  ({
    service,
    allServices,
    currency,
    isModification,
    type,
    contractStatus,
  }: ServiceItemProps) => {
    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const onChangeServiceStarted = (e: Date | null) => {
      if (!e) return;

      const checkExistingServiceStarted = (date: Date) => {
        return allServices?.some((service) =>
          DateTimeUtils.isSameDay(
            service?.tempValue?.serviceStarted,
            `${date}`,
          ),
        );
      };

      const findCurrentService = () => {
        if (isDraft) return null;

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
        return (
          currentService &&
          DateTimeUtils.isBefore(date.toString(), currentService.toString())
        );
      };

      const existingServiceStarted = checkExistingServiceStarted(e);
      const currentService = findCurrentService();
      const isBeforeCurrentService = checkIfBeforeCurrentService(
        e,
        currentService,
      );

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
        serviceStarted: e,
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
        (prev) => ({ ...prev, price: price ? parseFloat(price) : undefined }),
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
            size='xs'
            type='number'
            mask={Number}
            autofix={true}
            placeholder='0'
            className={inputClasses}
            onFocus={(e) => e.target.select()}
            min={type === 'one-time' ? -999999999999 : 0}
            value={service?.tempValue?.price?.toString() ?? ''}
            onChange={(e) => {
              updatePrice(e.target.value ?? '');
            }}
            onBlur={(e) =>
              !e.target.value?.length
                ? updatePrice('0')
                : updatePrice(e.target.value)
            }
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

        <Tooltip label='Service start date'>
          <div>
            <DatePickerUnderline2
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
      </div>
    );
  },
);
