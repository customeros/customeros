'use client';

import React, { useState, useEffect } from 'react';

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

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
function validateLocalPart(localPart: string) {
  if (localPart.startsWith('.') || localPart.endsWith('.')) {
    return 'The email address cannot start or end with a dot.';
  }

  const regex =
    /^(?:(?:[a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+)|(?:"(?:[\\"]|[^"\\])*")|(?:\.(?!\.)))+$/;

  if (!regex.test(localPart)) {
    if (localPart.includes('..')) {
      return 'The email address cannot contain consecutive dots.';
    }

    return 'The email address contains invalid characters.';
  }

  return '';
}

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
  bcc,
}: {
  formId: string;
  bcc?: string | null;
  email?: string | null;
  canPayWithCard?: boolean;
  invoicingEnabled?: boolean;
  canPayWithDirectDebitACH?: boolean;
  canPayWithDirectDebitSEPA?: boolean;
  canPayWithDirectDebitBacs?: boolean;
  setIsInvoiceProviderFocused: (newState: boolean) => void;
  setIsInvoiceProviderDetailsHovered: (newState: boolean) => void;
}) => {
  const [emailError, setEmailError] = useState<string | null>(null);
  useEffect(() => {
    if (email && email?.length > 0) {
      const emailValidation = validateLocalPart(email);
      setEmailError(emailValidation);
    }
  }, [email]);

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
          isInvalid={!!email?.length && Boolean(emailError)}
          error={emailError}
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
          isInvalid={!!bcc?.length && !emailRegex.test(bcc)}
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
