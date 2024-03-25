'use client';

import React, { useMemo } from 'react';

import { LogoUploader } from '@settings/components/LogoUploadComponent/LogoUploader';
import { PaymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/PaymentMethods';

import { FormSelect } from '@ui/form/SyncSelect';
import { Divider } from '@ui/presentation/Divider';
import { countryOptions } from '@shared/util/countryOptions';
import { FormMaskInput } from '@ui/form/Input/FormMaskInput';
import { FormInput, FormResizableInput } from '@ui/form/Input';
import { getCurrencyOptions } from '@shared/util/currencyOptions';

const VAT = {
  mask: 'AA 000 000 000',
  definitions: {
    A: /[A-Za-z]/,
    '0': /[0-9]/,
  },
  prepare: function (value: string, mask: { _value: string }) {
    if (mask._value.length < 2) {
      return value.toUpperCase();
    }

    return value;
  },
  format: function (value: string) {
    return value.toUpperCase();
  },
  parse: function (value: string) {
    return value.toUpperCase();
  },
};

export const TenantBillingPanelDetailsForm = ({
  setIsInvoiceProviderDetailsHovered,
  setIsInvoiceProviderFocused,
  formId,
  organizationName,
}: {
  formId: string;
  organizationName?: string | null;
  setIsInvoiceProviderFocused: (newState: boolean) => void;
  setIsInvoiceProviderDetailsHovered: (newState: boolean) => void;
}) => {
  const currencyOptions = useMemo(() => getCurrencyOptions(), []);

  return (
    <div className='w-full flex flex-col px-6 gap-4'>
      <LogoUploader />
      <FormInput
        autoComplete='off'
        label='Organization legal name'
        placeholder='Legal name'
        labelProps={{
          className: 'text-sm mb-0 font-semibold',
        }}
        name='legalName'
        formId={formId}
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
        onFocus={() => setIsInvoiceProviderFocused(true)}
        onBlur={() => setIsInvoiceProviderFocused(false)}
      />

      <FormSelect
        label='Base currency'
        placeholder='Invoice currency'
        isLabelVisible
        name='baseCurrency'
        formId={formId}
        options={currencyOptions ?? []}
      />

      <div
        className='flex flex-col'
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
      >
        <FormInput
          autoComplete='off'
          label='Billing address'
          placeholder='Address line 1'
          labelProps={{
            className: 'text-sm mb-0 font-semibold',
          }}
          name='addressLine1'
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
        <FormInput
          autoComplete='off'
          label='Billing address line 2'
          labelProps={{ style: { display: 'none' } }}
          name='addressLine2'
          placeholder='Address line 2'
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />

        <div className='flex space-x-2'>
          <FormInput
            autoComplete='off'
            label='Billing address locality'
            name='locality'
            labelProps={{ className: 'hidden' }}
            placeholder='City'
            formId={formId}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
          />
          <FormInput
            autoComplete='off'
            label='Billing address zip/Postal code'
            name='zip'
            placeholder='ZIP/Postal code'
            labelProps={{ className: 'hidden' }}
            formId={formId}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
          />
        </div>
        <FormSelect
          name='country'
          placeholder='Country'
          formId={formId}
          options={countryOptions}
        />

        <FormMaskInput
          options={{ opts: VAT }}
          formId={formId}
          name='vatNumber'
          autoComplete='off'
          label='VAT number'
          labelProps={{
            className: 'text-sm mb-0 font-semibold mt-4',
          }}
          className='overflow-ellipsis'
          placeholder='VAT number'
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
      </div>
      <div className='flex flex-col'>
        <div className='items-center relative'>
          <span className='tes-sm whitespace-nowrap mr-2 text-gray.500'>
            Email invoice
          </span>
          <Divider background='gray.200' />
        </div>

        <FormResizableInput
          formId={formId}
          autoComplete='off'
          label='From'
          labelProps={{
            className: 'text-sm mb-0 font-semibold inline-block pt-4',
          }}
          fontWeight='medium'
          name='sendInvoicesFrom'
          placeholder=''
          className='font-semibold'
          rightElement={'@invoices.customeros.com'}
          onFocus={() => setIsInvoiceProviderFocused(true)}
        />

        <FormInput
          autoComplete='off'
          label='BCC'
          labelProps={{
            className: 'text-sm mb-0 h-[100%] pt-4 font-semibold inline-block',
          }}
          formId={formId}
          name='sendInvoicesBcc'
          placeholder='BCC'
          type='email'
          className='overflow-ellipsis'
        />
      </div>
      <PaymentMethods formId={formId} organizationName={organizationName} />
    </div>
  );
};
