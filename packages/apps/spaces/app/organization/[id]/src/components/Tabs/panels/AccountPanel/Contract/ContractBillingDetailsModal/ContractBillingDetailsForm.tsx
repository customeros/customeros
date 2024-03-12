'use client';
import Link from 'next/link';
import React, { FC, useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';
import { Tooltip } from '@ui/overlay/Tooltip';
import { ModalBody } from '@ui/overlay/Modal';
import { FormUrlInput } from '@ui/form/UrlInput';
import { FormSelect } from '@ui/form/SyncSelect';
import { TenantBillingProfile } from '@graphql/types';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { countryOptions } from '@shared/util/countryOptions';
import { FormCheckbox } from '@ui/form/Checkbox/FormCheckbox';
import { getCurrencyOptions } from '@shared/util/currencyOptions';
import {
  Popover,
  PopoverBody,
  PopoverArrow,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

interface SubscriptionServiceModalProps {
  formId: string;
  currency?: string;
  isEmailValid: boolean;
  organizationName: string;
  hasNoBankAccounts: boolean;
  canAllowPayWithBankTransfer?: boolean;
  tenantBillingProfile?: TenantBillingProfile | null;
  onSetIsBillingDetailsHovered: (newState: boolean) => void;
  onSetIsBillingDetailsFocused: (newState: boolean) => void;
}

export const ContractBillingDetailsForm: FC<SubscriptionServiceModalProps> = ({
  formId,
  canAllowPayWithBankTransfer,
  isEmailValid,
  onSetIsBillingDetailsFocused,
  onSetIsBillingDetailsHovered,
  hasNoBankAccounts,
  currency,
  tenantBillingProfile,
  organizationName,
}) => {
  const currencyOptions = useMemo(() => getCurrencyOptions(), []);

  const tooltipContent = useMemo(() => {
    if (
      tenantBillingProfile?.canPayWithCard &&
      tenantBillingProfile?.canPayWithDirectDebitACH
    ) {
      return `If auto-payment fails, ${organizationName} can still pay using one of the other enabled payment options.`;
    }

    return '';
  }, [
    tenantBillingProfile?.canPayWithCard,
    tenantBillingProfile?.canPayWithDirectDebitACH,
    organizationName,
  ]);

  return (
    <ModalBody pb='0' gap={4} display='flex' flexDir='column' flex={1}>
      <FormUrlInput
        label='Link to contract'
        isLabelVisible
        labelProps={{
          fontSize: 'sm',
          mb: 0,
          fontWeight: 'semibold',
        }}
        formId={formId}
        name='contractUrl'
        textOverflow='ellipsis'
        placeholder='Paste or enter a contract link'
        autoComplete='off'
      />

      <FormInput
        label='Organization legal name'
        isLabelVisible
        labelProps={{
          fontSize: 'sm',
          mb: 0,
          fontWeight: 'semibold',
        }}
        onMouseEnter={() => onSetIsBillingDetailsHovered(true)}
        onMouseLeave={() => onSetIsBillingDetailsHovered(false)}
        onFocus={() => onSetIsBillingDetailsFocused(true)}
        onBlur={() => onSetIsBillingDetailsFocused(false)}
        formId={formId}
        name='organizationLegalName'
        textOverflow='ellipsis'
        placeholder='Organization legal name'
        autoComplete='off'
      />

      <Flex
        flexDir='column'
        onMouseEnter={() => onSetIsBillingDetailsHovered(true)}
        onMouseLeave={() => onSetIsBillingDetailsHovered(false)}
      >
        <FormInput
          label='Billing address'
          isLabelVisible
          labelProps={{
            fontSize: 'sm',
            mb: 0,
            fontWeight: 'semibold',
          }}
          formId={formId}
          name='addressLine1'
          textOverflow='ellipsis'
          placeholder='Address line 1'
          onFocus={() => onSetIsBillingDetailsFocused(true)}
          onBlur={() => onSetIsBillingDetailsFocused(false)}
          autoComplete='off'
        />
        <FormInput
          label='Address line 2'
          formId={formId}
          name='addressLine2'
          textOverflow='ellipsis'
          placeholder='Address line 2'
          onFocus={() => onSetIsBillingDetailsFocused(true)}
          onBlur={() => onSetIsBillingDetailsFocused(false)}
          autoComplete='off'
        />
        <Flex>
          <FormInput
            label='City'
            formId={formId}
            name='locality'
            textOverflow='ellipsis'
            placeholder='City'
            onFocus={() => onSetIsBillingDetailsFocused(true)}
            onBlur={() => onSetIsBillingDetailsFocused(false)}
            autoComplete='off'
          />
          <FormInput
            label='ZIP/Postal code'
            formId={formId}
            name='zip'
            textOverflow='ellipsis'
            placeholder='ZIP/Postal code'
            onFocus={() => onSetIsBillingDetailsFocused(true)}
            onBlur={() => onSetIsBillingDetailsFocused(false)}
            autoComplete='off'
          />
        </Flex>
        <FormSelect
          label='Country'
          placeholder='Country'
          name='country'
          formId={formId}
          options={countryOptions}
          onFocus={() => onSetIsBillingDetailsFocused(true)}
          onBlur={() => onSetIsBillingDetailsFocused(false)}
        />
      </Flex>
      <FormInput
        label='Send invoice to'
        isLabelVisible
        labelProps={{
          fontSize: 'sm',
          mb: 0,
          fontWeight: 'semibold',
        }}
        formId={formId}
        name='invoiceEmail'
        textOverflow='ellipsis'
        placeholder='Email'
        type='email'
        isInvalid={isEmailValid}
        onMouseEnter={() => onSetIsBillingDetailsHovered(true)}
        onMouseLeave={() => onSetIsBillingDetailsHovered(false)}
        onFocus={() => onSetIsBillingDetailsFocused(true)}
        onBlur={() => onSetIsBillingDetailsFocused(false)}
        autoComplete='off'
      />
      <FormSelect
        label='Billing currency'
        placeholder='Invoice currency'
        isLabelVisible
        name='currency'
        formId={formId}
        options={currencyOptions ?? []}
      />

      <Flex flexDirection='column'>
        <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap' mb={2}>
          Payment options
          {tooltipContent && (
            <Tooltip label={tooltipContent} shouldWrapChildren hasArrow>
              <InfoCircle boxSize={3} color='gray.400' ml={2} />
            </Tooltip>
          )}
        </Text>

        <FormSwitch
          name='payAutomatically'
          formId={formId}
          size='sm'
          label={
            <Text fontSize='sm' fontWeight='normal' whiteSpace='nowrap'>
              Auto-payment via Stripe
            </Text>
          }
        />
        <Flex flexDir='column' gap={2} ml={2}>
          <Tooltip
            shouldWrapChildren
            hasArrow
            label={
              tenantBillingProfile?.canPayWithCard
                ? ''
                : 'Credit or Debit card not enabled in Stripe'
            }
          >
            <FormCheckbox
              name='canPayWithCard'
              formId={formId}
              size='md'
              isInvalid={!tenantBillingProfile?.canPayWithCard}
            >
              <Text fontSize='sm' whiteSpace='nowrap'>
                Credit or Debit cards
              </Text>
            </FormCheckbox>
          </Tooltip>
          <Tooltip
            shouldWrapChildren
            hasArrow
            label={
              tenantBillingProfile?.canPayWithDirectDebitACH
                ? ''
                : 'Direct debit not enabled in Stripe'
            }
          >
            <FormCheckbox
              name='canPayWithDirectDebit'
              formId={formId}
              size='md'
              isInvalid={!tenantBillingProfile?.canPayWithDirectDebitACH}
            >
              <Text fontSize='sm' whiteSpace='nowrap'>
                Direct Debit via ACH
              </Text>
            </FormCheckbox>
          </Tooltip>
        </Flex>
        <FormSwitch
          name='payOnline'
          formId={formId}
          size='sm'
          label={
            <Text fontSize='sm' fontWeight='normal' whiteSpace='nowrap'>
              Online payment via Stripe
            </Text>
          }
        />
        <Popover placement='bottom-end' trigger='hover'>
          <PopoverTrigger>
            <FormSwitch
              name='canPayWithBankTransfer'
              isInvalid={
                !canAllowPayWithBankTransfer ||
                !tenantBillingProfile?.canPayWithBankTransfer
              }
              formId={formId}
              size='sm'
              label={
                <Text fontSize='sm' fontWeight='normal' whiteSpace='nowrap'>
                  Bank transfer
                </Text>
              }
            />
          </PopoverTrigger>
          <PopoverContent
            width='fit-content'
            bg='gray.700'
            color='white'
            mt={4}
            borderRadius='md'
            boxShadow='none'
            border='none'
          >
            <PopoverArrow bg='gray.700' />

            <PopoverBody display='flex'>
              <Text mr={2}>
                {!tenantBillingProfile?.canPayWithBankTransfer &&
                  'Bank transfer not enabled yet'}
                {tenantBillingProfile?.canPayWithBankTransfer &&
                canAllowPayWithBankTransfer &&
                hasNoBankAccounts
                  ? 'No bank accounts added yet'
                  : `None of your bank accounts hold ${currency}`}
              </Text>
              <Text as={Link} href='/settings?tab=billing' color='white'>
                Go to Settings
              </Text>
            </PopoverBody>
          </PopoverContent>
        </Popover>
      </Flex>
    </ModalBody>
  );
};
