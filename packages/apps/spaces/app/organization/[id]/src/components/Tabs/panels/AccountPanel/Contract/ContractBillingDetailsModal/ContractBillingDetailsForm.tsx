'use client';
import React, { FC, useMemo } from 'react';

import { useConnections } from '@integration-app/react';
import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';
import { useGetExternalSystemInstancesQuery } from '@settings/graphql/getExternalSystemInstances.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';
import { ModalBody } from '@ui/overlay/Modal';
import { Tooltip } from '@ui/overlay/Tooltip';
import { FormUrlInput } from '@ui/form/UrlInput';
import { FormSelect } from '@ui/form/SyncSelect';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { countryOptions } from '@shared/util/countryOptions';
import { FormCheckbox } from '@ui/form/Checkbox/FormCheckbox';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { getCurrencyOptions } from '@shared/util/currencyOptions';
import {
  BankAccount,
  ExternalSystemType,
  TenantBillingProfile,
} from '@graphql/types';
import { PaymentDetailsPopover } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/PaymentDetailsPopover';

interface SubscriptionServiceModalProps {
  formId: string;
  currency?: string;
  isEmailValid: boolean;
  organizationName: string;
  tenantBillingProfile?: TenantBillingProfile | null;
  bankAccounts: Array<BankAccount> | null | undefined;
  onSetIsBillingDetailsHovered: (newState: boolean) => void;
  onSetIsBillingDetailsFocused: (newState: boolean) => void;
}

export const ContractBillingDetailsForm: FC<SubscriptionServiceModalProps> = ({
  formId,
  isEmailValid,
  onSetIsBillingDetailsFocused,
  onSetIsBillingDetailsHovered,
  currency,
  tenantBillingProfile,
  organizationName,
  bankAccounts,
}) => {
  const client = getGraphQLClient();
  const { data: tenantSettingsData } = useTenantSettingsQuery(client);

  const { data } = useGetExternalSystemInstancesQuery(client);
  const currencyOptions = useMemo(() => getCurrencyOptions(), []);
  const availablePaymentMethodTypes = data?.externalSystemInstances.find(
    (e) => e.type === ExternalSystemType.Stripe,
  )?.stripeDetails?.paymentMethodTypes;
  const { items: iConnections } = useConnections();
  const isStripeActive = !!iConnections
    .map((item) => item.integration?.key)
    .find((e) => e === 'stripe');
  const tooltipContent = useMemo(() => {
    if (availablePaymentMethodTypes?.length && isStripeActive) {
      return `If auto-payment fails, ${organizationName} can still pay using one of the other enabled payment options.`;
    }

    return '';
  }, [isStripeActive, availablePaymentMethodTypes, organizationName]);

  const bankTransferPopoverContent = useMemo(() => {
    if (!tenantBillingProfile?.canPayWithBankTransfer) {
      return 'Bank transfer not enabled yet';
    }
    if (
      tenantBillingProfile?.canPayWithBankTransfer &&
      (!bankAccounts || bankAccounts.length === 0)
    ) {
      return 'No bank accounts added yet';
    }
    const accountIndexWithCurrency = bankAccounts?.findIndex(
      (account) => account.currency === currency,
    );

    if (accountIndexWithCurrency === -1 && currency) {
      return `None of your bank accounts hold ${currency}`;
    }
    if (!currency) {
      return `Please select contract currency to enable bank transfer`;
    }

    return '';
  }, [tenantBillingProfile, bankAccounts, currency]);

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

      {tenantSettingsData?.tenantSettings?.billingEnabled && (
        <>
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

          <Flex flexDirection='column' gap={2}>
            <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>
              Payment options
              {tooltipContent && (
                <Tooltip label={tooltipContent} hasArrow shouldWrapChildren>
                  <InfoCircle boxSize={3} color='gray.400' ml={2} />
                </Tooltip>
              )}
            </Text>

            <Flex flexDir='column' gap={2}>
              <PaymentDetailsPopover
                content={isStripeActive ? '' : 'No payment provider enabled'}
                withNavigation
              >
                <FormSwitch
                  name='payAutomatically'
                  formId={formId}
                  isInvalid={!isStripeActive}
                  size='sm'
                  labelProps={{ margin: 0 }}
                  label={
                    <Text fontSize='sm' fontWeight='normal' whiteSpace='nowrap'>
                      Auto-payment via Stripe
                    </Text>
                  }
                />
              </PaymentDetailsPopover>
              {isStripeActive && (
                <Flex flexDir='column' gap={2} ml={2}>
                  <Tooltip
                    label={
                      availablePaymentMethodTypes?.includes('card')
                        ? ''
                        : 'Credit or Debit card not enabled in Stripe'
                    }
                    placement='bottom-start'
                  >
                    <Box>
                      <FormCheckbox
                        name='canPayWithCard'
                        formId={formId}
                        size='md'
                        isInvalid={
                          !availablePaymentMethodTypes?.includes('card')
                        }
                      >
                        <Text fontSize='sm' whiteSpace='nowrap'>
                          Credit or Debit cards
                        </Text>
                      </FormCheckbox>
                    </Box>
                  </Tooltip>
                  <Tooltip
                    label={
                      availablePaymentMethodTypes?.includes('bacs_debit')
                        ? ''
                        : 'Direct debit not enabled in Stripe'
                    }
                    placement='bottom-start'
                  >
                    <Box>
                      <FormCheckbox
                        name='canPayWithDirectDebit'
                        formId={formId}
                        size='md'
                        isInvalid={
                          !availablePaymentMethodTypes?.includes('bacs_debit')
                        }
                      >
                        <Text fontSize='sm' whiteSpace='nowrap'>
                          Direct Debit via ACH
                        </Text>
                      </FormCheckbox>
                    </Box>
                  </Tooltip>
                </Flex>
              )}
            </Flex>

            <PaymentDetailsPopover
              content={isStripeActive ? '' : 'No payment provider enabled'}
              withNavigation
            >
              <FormSwitch
                name='payOnline'
                formId={formId}
                isInvalid={!isStripeActive}
                size='sm'
                labelProps={{
                  margin: 0,
                }}
                label={
                  <Text fontSize='sm' fontWeight='normal' whiteSpace='nowrap'>
                    Online payment via Stripe
                  </Text>
                }
              />
            </PaymentDetailsPopover>

            <PaymentDetailsPopover
              withNavigation
              content={bankTransferPopoverContent}
            >
              <FormSwitch
                name='canPayWithBankTransfer'
                isInvalid={!!bankTransferPopoverContent.length}
                formId={formId}
                size='sm'
                labelProps={{
                  margin: 0,
                }}
                label={
                  <Text fontSize='sm' fontWeight='normal' whiteSpace='nowrap'>
                    Bank transfer
                  </Text>
                }
              />
            </PaymentDetailsPopover>
          </Flex>
        </>
      )}
    </ModalBody>
  );
};
