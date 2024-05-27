import { useRef, useState, useEffect } from 'react';

import { toZonedTime } from 'date-fns-tz';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input/Input';
import { Delete } from '@ui/media/icons/Delete';
import { DateTimeUtils } from '@spaces/utils/date';
import { MaskInput } from '@ui/form/Input/MaskInput';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { FlipBackward } from '@ui/media/icons/FlipBackward';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { NumberInput } from '@ui/form/NumberInput/NumberInput';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { billedTypeOptions } from '@organization/components/Tabs/panels/AccountPanel/utils';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';
import {
  Select,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select';
import { Currency } from '@organization/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/Currency';
import { BulkUpdateServiceLineItem } from '@organization/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/ServiceLineItemsModal.dto';

import { ServiceLineItemInputWrapper } from './ServiceLineItemInputWrapper';

type DateInputValue = null | string | number | Date;

const [_, _1, ...subscriptionOptions] = billedTypeOptions;
interface ServiceLineItemProps {
  index: number;
  currency?: string | null;
  service: BulkUpdateServiceLineItem;
  prevServiceLineItemData?: ServiceLineItem;
  onChange: (updatedService: BulkUpdateServiceLineItem) => void;
}

export const ServiceLineItemRow = ({
  service,
  onChange,
  index,
  currency,
}: ServiceLineItemProps) => {
  const handleChange = (
    field: keyof BulkUpdateServiceLineItem,
    value: string | boolean | Date,
  ) => {
    onChange({ ...service, [field]: value });
  };
  const [isOpen, setIsOpen] = useState(false);

  const nameRef = useRef<HTMLInputElement | null>(null);
  const currencySymbol =
    formatCurrency(0, 0, currency || 'USD')?.split('0')[0] ?? '$';
  const handleTypeChange = (newValue: string) => {
    if (newValue === 'RECURRING') {
      onChange({
        ...service,
        type: 'RECURRING',
        billed: BilledType.Monthly,
      });

      return;
    }
    onChange({
      ...service,
      type: newValue,
      billed: newValue as BilledType,
    });
  };

  useEffect(() => {
    if (service.name === 'Unnamed' && nameRef?.current) {
      nameRef.current?.focus();
      nameRef.current?.setSelectionRange(0, service.name.length);
    }
  }, [service.name, nameRef]);

  const handleDateInputChange = (data?: DateInputValue) => {
    if (!data) return handleChange('serviceStarted', false);
    const date = new Date(data);

    const normalizedDate = new Date(
      Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()),
    );
    handleChange('serviceStarted', normalizedDate);
  };

  return (
    <div className='flex justify-between items-center gap-1 border-b border-gray-300 relative pr-4 [&_.delete-button]:opacity-0 [&_.delete-button]:transition-opacity duration-150 ease-in [&_.delete-button]:hover:opacity-100 '>
      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Input
          name='name'
          aria-label='Name'
          placeholder='Name'
          autoFocus
          value={service.name}
          ref={nameRef}
          onChange={(event) => handleChange('name', event.target.value)}
          key={`name-${index}`}
          autoComplete='off'
          className={cn(service.isDeleted ? 'line-through' : '')}
        />
      </ServiceLineItemInputWrapper>
      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Select
          aria-label='Type'
          placeholder='Type'
          name='type'
          value={
            [
              BilledType.Quarterly,
              BilledType.Monthly,
              BilledType.Annually,
            ].includes(service.billed as BilledType)
              ? TypeOptions[0]
              : TypeOptions.find((e) => e.value === service.billed)
          }
          onChange={(newValue) => {
            handleTypeChange(newValue.value);
          }}
          options={TypeOptions}
          classNames={{
            container: () =>
              getContainerClassNames(cn({ 'line-through': service.isDeleted })),
          }}
        />
      </ServiceLineItemInputWrapper>
      <ServiceLineItemInputWrapper width='10%' isDeleted={service.isDeleted}>
        <NumberInput
          placeholder='10'
          aria-label='Quantity'
          min={1}
          name='quantity'
          value={service.quantity}
          autoComplete='off'
          className={cn(service.isDeleted ? 'line-through' : '')}
          onChange={(event) => handleChange('quantity', event.target.value)}
        />
      </ServiceLineItemInputWrapper>
      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Currency
          name='price'
          className={cn(
            service.isDeleted ? 'line-through' : '',
            'w-full text-sm',
          )}
          placeholder='Per license'
          label='Price/qty'
          labelProps={{ className: 'hidden' }}
          value={`${service.price}`}
          currency={currencySymbol}
          onValueChange={(value) => {
            handleChange('price', value);
          }}
        />
      </ServiceLineItemInputWrapper>
      <ServiceLineItemInputWrapper width='10%' isDeleted={service.isDeleted}>
        {service.type === 'RECURRING' ? (
          <Select
            aria-label='Recurring'
            placeholder='Frequency'
            name='billed'
            onChange={(newValue) => handleChange('billed', newValue.value)}
            options={subscriptionOptions}
            value={subscriptionOptions.find((e) => e.value === service.billed)}
            classNames={{
              menuList: () => getMenuListClassNames('min-w-[100px]'),
              container: () =>
                getContainerClassNames(
                  cn({ 'line-through': service.isDeleted }),
                ),
            }}
          />
        ) : (
          <p className='text-gray-400' color='gray.400'>
            N/A
          </p>
        )}
      </ServiceLineItemInputWrapper>

      <ServiceLineItemInputWrapper width='10%' isDeleted={service.isDeleted}>
        <MaskInput
          placeholder='0'
          aria-label='VAT'
          min={0}
          name='vatRate'
          autoComplete='off'
          value={`${service.vatRate}`}
          symbol='%'
          onValueChange={(value) => {
            handleChange('vatRate', value);
          }}
          className={cn(service.isDeleted ? 'line-through' : '')}
        />
      </ServiceLineItemInputWrapper>

      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
          <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
            <span
              className={cn(
                service.isDeleted ? 'line-through' : '',
                service.serviceStarted ? 'text-gray-700' : 'text-gray-400',
                'cursor-pointer whitespace-pre pb-[1px] text-base border-t-[1px] border-transparent hover:text-gray-700',
              )}
            >
              {service.serviceStarted
                ? DateTimeUtils.format(
                    (new Date(service.serviceStarted) as Date)?.toISOString(),
                    DateTimeUtils.dateWithAbreviatedMonth,
                  )
                : 'Start date'}
            </span>
          </PopoverTrigger>
          <PopoverContent
            align='center'
            side='bottom'
            sticky='always'
            onOpenAutoFocus={(el) => el.preventDefault()}
            onClick={(e) => e.stopPropagation()}
          >
            <DatePicker
              formId='service-line-template-date-picker'
              name='startDate'
              defaultValue={
                service.serviceStarted
                  ? toZonedTime(service.serviceStarted, 'UTC')
                  : null
              }
              onChange={(date) => {
                handleDateInputChange(date as Date);
              }}
            />
          </PopoverContent>
        </Popover>
      </ServiceLineItemInputWrapper>
      <IconButton
        aria-label='Delete'
        className='delete-button absolute right-[-4px]'
        icon={
          service.isDeleted ? (
            <FlipBackward className='text-gray-400' color='gray.400' />
          ) : (
            <Delete className='text-gray-400' color='gray.400' />
          )
        }
        variant='ghost'
        size='xs'
        onClick={() => handleChange('isDeleted', !service.isDeleted)}
      />
    </div>
  );
};
export const TypeOptions: SelectOption<string>[] = [
  { label: 'Recurring', value: 'RECURRING' },
  { label: 'Per-use', value: BilledType.Usage },
  { label: 'One-time', value: BilledType.Once },
];
