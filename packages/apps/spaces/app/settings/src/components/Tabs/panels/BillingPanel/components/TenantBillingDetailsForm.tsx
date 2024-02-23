'use client';

import React from 'react';

import { LogoUploader } from '@settings/components/LogoUploadComponent/LogoUploader';
import { VatInput } from '@settings/components/Tabs/panels/BillingPanel/components/VatInput';
import { PaymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/PaymentMethods';

import { Flex } from '@ui/layout/Flex';
import { CardBody } from '@ui/layout/Card';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { countryOptions } from '@shared/util/countryOptions';

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const TenantBillingPanelDetailsForm = ({
  setIsInvoiceProviderDetailsHovered,
  setIsInvoiceProviderFocused,
  formId,
  invoicingEnabled,
  canPayWithCard,
  canPayWithDirectDebitACH,
  canPayWithDirectDebitSEPA,
  canPayWithDirectDebitBacs,
  email,
}: {
  formId: string;
  email?: string | null;
  canPayWithCard?: boolean;
  invoicingEnabled?: boolean;
  canPayWithDirectDebitACH?: boolean;
  canPayWithDirectDebitSEPA?: boolean;
  canPayWithDirectDebitBacs?: boolean;
  setIsInvoiceProviderFocused: (newState: boolean) => void;
  setIsInvoiceProviderDetailsHovered: (newState: boolean) => void;
}) => {
  return (
    <CardBody
      as={Flex}
      flexDir='column'
      px='6'
      w='full'
      gap={4}
      opacity={invoicingEnabled ? 1 : 0}
    >
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

        <FormInput
          autoComplete='off'
          label='Send invoice from'
          isLabelVisible
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            mt: 4,
            fontWeight: 'semibold',
          }}
          formId={formId}
          name='sendInvoicesFrom'
          textOverflow='ellipsis'
          placeholder='Email'
          type='email'
          isInvalid={!!email?.length && !emailRegex.test(email)}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
      </Flex>

      <PaymentMethods
        formId={formId}
        canPayWithCard={canPayWithCard}
        canPayWithDirectDebitACH={canPayWithDirectDebitACH}
        canPayWithDirectDebitSEPA={canPayWithDirectDebitSEPA}
        canPayWithDirectDebitBacs={canPayWithDirectDebitBacs}
      />
      {/*<Flex justifyContent='space-between' alignItems='center'>*/}
      {/*  <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>*/}
      {/*    Bank transfer*/}
      {/*  </Text>*/}
      {/*  <Switch size='sm' />*/}
      {/*</Flex>*/}
    </CardBody>
  );
};
