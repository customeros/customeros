'use client';
import React, { useRef, useEffect } from 'react';
import { DatePicker as ReactDatePicker } from 'react-date-picker';

import { utcToZonedTime } from 'date-fns-tz';

import { Flex } from '@ui/layout/Flex';
import { Input } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Select } from '@ui/form/SyncSelect';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';
import { DateTimeUtils } from '@spaces/utils/date';
import { SelectOption } from '@shared/types/SelectOptions';
import { FlipBackward } from '@ui/media/icons/FlipBackward';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { NumberInput, NumberInputField } from '@ui/form/NumberInput';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { Currency } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/Currency';
import { BulkUpdateServiceLineItem } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/ServiceLineItemsModal.dto';

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
  prevServiceLineItemData,
}: ServiceLineItemProps) => {
  const handleChange = (
    field: keyof BulkUpdateServiceLineItem,
    value: string | boolean | Date,
  ) => {
    onChange({ ...service, [field]: value });
  };
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
    <Flex
      justifyContent='space-between'
      alignItems='center'
      gap={1}
      position='relative'
      pr='20px'
      borderBottom='1px solid'
      borderColor='gray.300'
      sx={{
        '.delete-button': {
          opacity: '0',
          transition: 'opacity 0.15s ease-in',
        },
        '&:hover .delete-button': {
          opacity: '1',
        },
      }}
    >
      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Input
          name='name'
          aria-label='Name'
          fontSize='sm'
          placeholder='Name'
          autoFocus
          value={service.name}
          ref={nameRef}
          onChange={(event) => handleChange('name', event.target.value)}
          key={`name-${index}`}
          textDecoration={service.isDeleted ? 'line-through' : 'unset'}
          autoComplete='off'
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
          chakraStyles={{
            container: (props, state) => {
              return {
                minHeight: 'unset',
                borderBottom: 'none',

                '& *': {
                  fontSize: 'sm',
                  textDecoration: service.isDeleted ? 'line-through' : 'unset',
                },
              };
            },
            input: (props, state) => {
              return { minHeight: 'auto', fontSize: 'sm' };
            },
            option: (props, state) => {
              return { fontSize: 'sm' };
            },
            menuList: (props, state) => ({
              minW: '130px',
            }),
          }}
        />
      </ServiceLineItemInputWrapper>
      <ServiceLineItemInputWrapper width='10%' isDeleted={service.isDeleted}>
        <NumberInput value={service.quantity}>
          <NumberInputField
            placeholder='10'
            aria-label='Quantity'
            textDecoration={service.isDeleted ? 'line-through' : 'unset'}
            min={1}
            name='quantity'
            fontSize='sm'
            value={service.quantity}
            p={0}
            autoComplete='off'
            onChange={(event) => handleChange('quantity', event.target.value)}
          />
        </NumberInput>
      </ServiceLineItemInputWrapper>
      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Currency
          name='price'
          w='full'
          placeholder='Per license'
          label='Price/qty'
          value={`${service.price}`}
          fontSize='sm'
          sx={{
            '&': {
              fontSize: '14px !important',
              textDecoration: service.isDeleted ? 'line-through' : 'unset',
            },
          }}
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
            chakraStyles={{
              container: (props, state) => {
                return {
                  minHeight: 'unset',
                  fontSize: 'sm',
                  borderBottom: 'none',
                  '& *': {
                    fontSize: 'sm',
                    textDecoration: service.isDeleted
                      ? 'line-through'
                      : 'unset',
                  },
                };
              },
              menuList: (props, state) => ({
                minW: '100px',
              }),
              input: (props, state) => {
                return {
                  minHeight: 'auto',
                  fontSize: 'sm',
                };
              },
              option: (props, state) => {
                return { fontSize: 'sm' };
              },
            }}
          />
        ) : (
          <Text color='gray.400'>N/A</Text>
        )}
      </ServiceLineItemInputWrapper>

      <ServiceLineItemInputWrapper width='10%' isDeleted={service.isDeleted}>
        <NumberInput value={`${service.vatRate}%`} display='flex'>
          <NumberInputField
            placeholder='0'
            aria-label='VAT'
            textDecoration={service.isDeleted ? 'line-through' : 'unset'}
            min={0}
            name='vatRate'
            fontSize='sm'
            p={0}
            onChange={(event) => {
              handleChange('vatRate', event.target.value.replace('%', ''));
            }}
          />
        </NumberInput>
      </ServiceLineItemInputWrapper>

      <ServiceLineItemInputWrapper width='15%' isDeleted={service.isDeleted}>
        <Flex
          sx={{
            '& .react-date-picker__calendar-button': {
              pl: 0,
            },
            '& .react-date-picker__calendar': {
              inset: `${'120% 0px auto auto'} !important`,
            },
            '& .react-date-picker__clear-button': {
              top: '7px',
            },
            '& .react-calendar__month-view__weekdays__weekday': {
              textTransform: 'capitalize',
            },
          }}
        >
          <ReactDatePicker
            id='service-line-template-date-picker'
            name='startDate'
            clearIcon={null}
            onChange={(event) => handleDateInputChange(event as DateInputValue)}
            defaultValue={
              service.serviceStarted
                ? utcToZonedTime(service.serviceStarted, 'UTC')
                : null
            }
            formatShortWeekday={(_, date) =>
              DateTimeUtils.format(
                date.toISOString(),
                DateTimeUtils.shortWeekday,
              )
            }
            formatMonth={(_, date) =>
              DateTimeUtils.format(
                date.toISOString(),
                DateTimeUtils.abreviatedMonth,
              )
            }
            calendarIcon={
              <Flex alignItems='center'>
                <Text
                  color={service.serviceStarted ? 'gray.700' : 'gray.400'}
                  role='button'
                  textDecoration={service.isDeleted ? 'line-through' : 'unset'}
                >
                  {service.serviceStarted
                    ? DateTimeUtils.format(
                        (
                          new Date(service.serviceStarted) as Date
                        )?.toISOString(),
                        DateTimeUtils.dateWithAbreviatedMonth,
                      )
                    : 'Start date'}
                </Text>
              </Flex>
            }
          />
        </Flex>
      </ServiceLineItemInputWrapper>
      <IconButton
        position='absolute'
        aria-label='Delete'
        className='delete-button'
        icon={
          service.isDeleted ? (
            <FlipBackward className='text-gray-400' />
          ) : (
            <Delete className='text-gray-400' />
          )
        }
        variant='ghost'
        size='xs'
        right={-1}
        onClick={() => handleChange('isDeleted', !service.isDeleted)}
      />
    </Flex>
  );
};
export const TypeOptions: SelectOption<string>[] = [
  { label: 'Recurring', value: 'RECURRING' },
  { label: 'Per-use', value: BilledType.Usage },
  { label: 'One-time', value: BilledType.Once },
];
