'use client';
import React, { useRef, useEffect } from 'react';

import { Currency } from '@settings/components/Tabs/panels/BillingPanel/Currency';

import { Flex } from '@ui/layout/Flex';
import { Input } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { BilledType } from '@graphql/types';
import { Select } from '@ui/form/SyncSelect';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';
import { SelectOption } from '@shared/types/SelectOptions';
import { FlipBackward } from '@ui/media/icons/FlipBackward';
import { NumberInput, NumberInputField } from '@ui/form/NumberInput';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

import { ServiceItem } from './type';
import { ServiceLineItemInputWrapper } from './ServiceLineItemInputWrapper';

const [_, _1, ...subscriptionOptions] = billedTypeOptions;
interface ServiceLineItemProps {
  index: number;
  service: ServiceItem;
  currency?: string | null;
  onChange: (updatedService: ServiceItem) => void;
}

export const ServiceLineItemRow = ({
  service,
  onChange,
  index,
  currency,
}: ServiceLineItemProps) => {
  const handleChange = (field: keyof ServiceItem, value: string | boolean) => {
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
        button: {
          opacity: '0',
          transition: 'opacity 0.15s ease-in',
        },
        '&:hover button': {
          opacity: '1',
        },
      }}
    >
      <ServiceLineItemInputWrapper width='30%' isDeleted={service.isDeleted}>
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
      <ServiceLineItemInputWrapper width='20%' isDeleted={service.isDeleted}>
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
      <ServiceLineItemInputWrapper width='20%' isDeleted={service.isDeleted}>
        <Currency
          name='price'
          w='full'
          placeholder='Per license'
          label='Price/qty'
          step={0.01}
          min={0.01}
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
      <ServiceLineItemInputWrapper width='20%' isDeleted={service.isDeleted}>
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

      <IconButton
        position='absolute'
        aria-label='Delete'
        icon={
          service.isDeleted ? (
            <FlipBackward color='gray.400' />
          ) : (
            <Delete color='gray.400' />
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
