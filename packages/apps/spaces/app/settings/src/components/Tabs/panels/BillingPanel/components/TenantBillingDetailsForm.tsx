'use client';

import React, { useMemo } from 'react';

import { LogoUploader } from '@settings/components/LogoUploadComponent/LogoUploader';
import { VatInput } from '@settings/components/Tabs/panels/BillingPanel/components/VatInput';
import { PaymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/PaymentMethods';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { CardBody } from '@ui/layout/Card';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { Divider } from '@ui/presentation/Divider';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { SelectOption } from '@shared/types/SelectOptions';
import { countryOptions } from '@shared/util/countryOptions';
import { getCurrencyOptions } from '@shared/util/currencyOptions';

export const TenantBillingPanelDetailsForm = ({
  setIsInvoiceProviderDetailsHovered,
  setIsInvoiceProviderFocused,
  formId,
  organizationName,
  check,
  sendInvoicesFrom,
  country,
}: {
  formId: string;
  check?: boolean | null;
  sendInvoicesFrom?: string;
  organizationName?: string | null;
  country?: SelectOption<string> | null;
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
        <FormSelect
          name='country'
          placeholder='Country'
          label='Country'
          isLabelVisible
          formId={formId}
          options={countryOptions}
        />
        <FormInput
          autoComplete='off'
          label='Billing address'
          placeholder='Address line 1'
          isLabelVisible
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            mt: 2,
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
        <FormInput
          autoComplete='off'
          label='Billing address line 2'
          name='addressLine2'
          placeholder='Address line 2'
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
        {country?.value === 'US' && (
          <FormInput
            label='State'
            name='region'
            placeholder='State'
            formId={formId}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
          />
        )}

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

        {sendInvoicesFrom && (
          <Tooltip label='This email is configured by CustomerOs'>
            <FormInput
              formId={formId}
              autoComplete='off'
              cursor='not-allowed'
              label='From'
              labelProps={{
                fontSize: 'sm',
                mb: 0,
                mt: 4,
                fontWeight: 'semibold',
              }}
              isLabelVisible
              isReadOnly
              name='sendInvoicesFrom'
              placeholder=''
              onFocus={() => setIsInvoiceProviderFocused(true)}
            />
          </Tooltip>
        )}

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
      <Box>
        <FormSwitch
          size='sm'
          name='check'
          formId={formId}
          label='Checks'
          fontWeight='semibold'
          labelProps={{
            fontSize: 'sm',
            fontWeight: 'semibold',
            margin: 0,
          }}
        />
        {check && <Text>Want to pay by check? Contact {sendInvoicesFrom}</Text>}
      </Box>
    </CardBody>
  );
};
