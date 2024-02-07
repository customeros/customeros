'use client';
import React, { FC, useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { ModalBody } from '@ui/overlay/Modal';
import { FormUrlInput } from '@ui/form/UrlInput';
import { FormSelect } from '@ui/form/SyncSelect';
import { countryOptions } from '@shared/util/countryOptions';
import { getCurrencyOptions } from '@shared/util/currencyOptions';

interface SubscriptionServiceModalProps {
  formId: string;
  isEmailValid: boolean;
  onSetIsBillingDetailsHovered: (newState: boolean) => void;
  onSetIsBillingDetailsFocused: (newState: boolean) => void;
}

export const ContractBillingDetailsForm: FC<SubscriptionServiceModalProps> = ({
  formId,
  isEmailValid,
  onSetIsBillingDetailsFocused,
  onSetIsBillingDetailsHovered,
}) => {
  const currencyOptions = useMemo(() => getCurrencyOptions(), []);

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
        label='Currency'
        placeholder='Invoice currency'
        isLabelVisible
        name='currency'
        formId={formId}
        options={currencyOptions ?? []}
      />
    </ModalBody>
  );
};
