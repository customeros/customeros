'use client';

import React, { useMemo } from 'react';

import { LogoUploader } from '@settings/components/LogoUploadComponent/LogoUploader';
import { VatInput } from '@settings/components/Tabs/panels/BillingPanel/components/VatInput';
import { PaymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/PaymentMethods';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { CardBody } from '@ui/layout/Card';
import { FormSelect } from '@ui/form/SyncSelect';
import { Divider } from '@ui/presentation/Divider';
import { countryOptions } from '@shared/util/countryOptions';
import { FormInput, FormResizableInput } from '@ui/form/Input';
import { getCurrencyOptions } from '@shared/util/currencyOptions';

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
    <CardBody as={Flex} flexDir='column' px='6' w='full' gap={4}>
      <LogoUploader />
      <FormInput
        autoComplete='off'
        label='Organization legal name'
        placeholder='Legal name'
        isLabelVisible
        labelProps={{
          fontSize: 'sm',
          mb: 0,
          fontWeight: 'semibold',
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

      <Flex
        flexDir='column'
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
      >
        <FormInput
          autoComplete='off'
          label='Billing address'
          placeholder='Address line 1'
          isLabelVisible
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            fontWeight: 'semibold',
          }}
          name='addressLine1'
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
        <FormInput
          autoComplete='off'
          label='Billing address line 2'
          name='addressLine2'
          placeholder='Address line 2'
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />

        <Flex gap={2}>
          <FormInput
            autoComplete='off'
            label='Billing address locality'
            name='locality'
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
            formId={formId}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
          />
        </Flex>
        <FormSelect
          name='country'
          placeholder='Country'
          formId={formId}
          options={countryOptions}
        />
        <VatInput
          formId={formId}
          name='vatNumber'
          autoComplete='off'
          label='VAT number'
          isLabelVisible
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            mt: 4,
            fontWeight: 'semibold',
          }}
          textOverflow='ellipsis'
          placeholder='VAT number'
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
      </Flex>
      <Flex flexDir='column'>
        <Flex position='relative' alignItems='center'>
          <Text fontSize='sm' whiteSpace='nowrap' mr={2} color='gray.500'>
            Email invoice
          </Text>
          <Divider background='gray.200' />
        </Flex>

        <FormResizableInput
          formId={formId}
          autoComplete='off'
          label='From'
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            mt: 4,
            fontWeight: 'semibold',
          }}
          fontWeight='medium'
          isLabelVisible
          name='sendInvoicesFrom'
          placeholder=''
          rightElement={'@invoices.customeros.com'}
          onFocus={() => setIsInvoiceProviderFocused(true)}
        />

        <FormInput
          autoComplete='off'
          label='BCC'
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            mt: 4,
            fontWeight: 'semibold',
          }}
          isLabelVisible
          formId={formId}
          name='sendInvoicesBcc'
          textOverflow='ellipsis'
          placeholder='BCC'
          type='email'
        />
      </Flex>
      <PaymentMethods formId={formId} organizationName={organizationName} />
    </CardBody>
  );
};
