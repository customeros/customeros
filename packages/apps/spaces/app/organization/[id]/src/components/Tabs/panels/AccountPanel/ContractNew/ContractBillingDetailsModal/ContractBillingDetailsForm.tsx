'use client';
import React, { FC, useMemo } from 'react';
import { useField } from 'react-inverted-form';

import { useConnections } from '@integration-app/react';
import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';
import { useGetExternalSystemInstancesQuery } from '@settings/graphql/getExternalSystemInstances.generated';

import { Button } from '@ui/form/Button/Button';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { FormSelect } from '@ui/form/Select/FormSelect';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { getMenuListClassNames } from '@ui/form/Select';
import { Divider } from '@ui/presentation/Divider/Divider';
import { FormCheckbox } from '@ui/form/Checkbox/FormCheckbox';
import { currencyOptions } from '@shared/util/currencyOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import {
  Currency,
  BankAccount,
  ExternalSystemType,
  TenantBillingProfile,
} from '@graphql/types';
import {
  paymentDueOptions,
  contractBillingCycleOptions,
} from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { CommittedPeriodInput } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/CommittedPeriodInput';
import { PaymentDetailsPopover } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/PaymentDetailsPopover';

import { ContractUploader } from './ContractUploader';

interface SubscriptionServiceModalProps {
  formId: string;
  currency?: string;
  contractId: string;
  payAutomatically?: boolean | null;
  tenantBillingProfile?: TenantBillingProfile | null;
  bankAccounts: Array<BankAccount> | null | undefined;
}

export const ContractBillingDetailsForm: FC<SubscriptionServiceModalProps> = ({
  formId,
  contractId,
  currency,
  tenantBillingProfile,
  bankAccounts,
  payAutomatically,
}) => {
  const client = getGraphQLClient();
  const { data: tenantSettingsData } = useTenantSettingsQuery(client);
  const { data } = useGetExternalSystemInstancesQuery(client);
  const availablePaymentMethodTypes = data?.externalSystemInstances.find(
    (e) => e.type === ExternalSystemType.Stripe,
  )?.stripeDetails?.paymentMethodTypes;
  const { items: iConnections } = useConnections();
  const isStripeActive = !!iConnections
    .map((item) => item.integration?.key)
    .find((e) => e === 'stripe');

  const { getInputProps } = useField('autoRenew', formId);
  const { onChange: onChangeAutoRenew, value: autoRenewValue } =
    getInputProps();
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

  const paymentMethod = useMemo(() => {
    let method;
    switch (currency) {
      case Currency.Gbp:
        method = 'Bacs';
        break;
      case Currency.Usd:
        method = 'ACH';
        break;
      default:
        method = 'SEPA';
    }

    return method;
  }, [currency]);

  return (
    <ModalBody className='flex flex-col flex-1 p-0'>
      <ul className='mb-2 list-disc ml-5'>
        <li className='text-sm '>
          <div className='flex items-baseline'>
            <CommittedPeriodInput formId={formId} />

            <span className='whitespace-nowrap mr-1'>contract, starting </span>

            <DatePickerUnderline formId={formId} name='serviceStarted' />
          </div>
        </li>
        <li className='text-sm mt-1.5'>
          <div className='flex items-baseline'>
            Live until 2 Aug 2024,{' '}
            <Button
              variant='ghost'
              size='sm'
              className='font-normal text-sm p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent after:content-[""] after:h-[1px] after:w-[100%] after:absolute after:bottom-[1px] after:left-0 after:bg-gray-500 hover:after:bg-gray-500 focus:after:bg-gray-700'
              onClick={() => onChangeAutoRenew(!autoRenewValue)}
            >
              {autoRenewValue ? 'auto-renews' : 'non auto-renewing'}
            </Button>
          </div>
        </li>
        <li className='text-sm '>
          <div className='flex items-baseline'>
            <span className='whitespace-nowrap'>Contracting in</span>
            <div>
              <FormSelect
                className='text-sm inline min-h-1 max-h-3 border-none hover:border-none focus:border-none w-fit ml-1 mt-0 after:content-[""] after:h-[0.5px] after:w-[1.7rem] after:absolute after:bottom-[4px] after:left-0 after:bg-gray-500 hover:after:bg-gray-500 focus:after:bg-gray-700'
                label='Currency'
                classNames={{
                  menuList: () => getMenuListClassNames('min-w-[120px]'),
                }}
                placeholder='Invoice currency'
                name='currency'
                formId={formId}
                options={currencyOptions ?? []}
                size='xs'
              />
            </div>
          </div>
        </li>
      </ul>

      {tenantSettingsData?.tenantSettings?.billingEnabled && (
        <>
          <div className='flex relative items-center h-8 mb-1'>
            <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
              Billing policy
            </p>
            <Divider />
          </div>
          <ul className='mb-2 list-disc ml-5'>
            <li className='text-sm '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap mr-1'>Billing starts </span>

                <DatePickerUnderline formId={formId} name='invoicingStarted' />
              </div>
            </li>
            <li className='text-sm '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap mr-1'>
                  Invoices are sent
                </span>
                <div>
                  <FormSelect
                    className='text-sm inline min-h-1 max-h-5 min-w-[50px] border-none hover:border-none focus:border-none w-fit mt-0 ml-0 mr-1 after:content-[""] after:h-[0.5px] after:w-[95%] after:m-r-1 after:absolute after:bottom-[4px] after:left-0 after:bg-gray-500 hover:after:bg-gray-500 focus:after:bg-gray-700'
                    label='billing period'
                    placeholder='billing period'
                    name='billingCycle'
                    formId={formId}
                    options={contractBillingCycleOptions}
                    size='xs'
                    classNames={{
                      menuList: () => getMenuListClassNames('min-w-[120px]'),
                    }}
                  />
                </div>
                <span className='whitespace-nowrap mr-1'>
                  on the billing start day
                </span>
              </div>
            </li>
            <li className='text-sm '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap mr-1'>Customer has</span>
                <div>
                  <FormSelect
                    className='text-sm inline min-h-1 max-h-5 border-none hover:border-none focus:border-none w-fit mt-0 ml-0 mr-1  after:content-[""] after:h-[0.5px] after:w-[95%] after:m-r-1 after:absolute after:bottom-[4px] after:left-0 after:bg-gray-500 hover:after:bg-gray-500 focus:after:bg-gray-700'
                    label='Payment due'
                    placeholder='0 days'
                    name='dueDays'
                    formId={formId}
                    options={paymentDueOptions}
                    formatOptionLabel={(option, formatOptionLabelMeta) => {
                      if (formatOptionLabelMeta.context === 'value') {
                        return `${option.value} days`;
                      }

                      return option.label;
                    }}
                    classNames={{
                      menuList: () => getMenuListClassNames('min-w-[120px]'),
                    }}
                    size='xs'
                  />
                </div>
                <span className='whitespace-nowrap mr-1'>to pay</span>
              </div>
            </li>
          </ul>
          <div className='flex relative items-center h-8 '>
            <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
              Payment options
            </p>
            <Divider />
          </div>

          <div className='flex flex-col gap-1 mb-2'>
            <div className='flex flex-col gap-1'>
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
                    <div className='text-sm font-normal whitespace-nowrap'>
                      Auto-payment via Stripe
                    </div>
                  }
                />
              </PaymentDetailsPopover>
              {isStripeActive && payAutomatically && (
                <div className='flex flex-col gap-1 ml-2'>
                  <Tooltip
                    label={
                      availablePaymentMethodTypes?.includes('card')
                        ? ''
                        : 'Credit or Debit card not enabled in Stripe'
                    }
                    side='bottom'
                    align='end'
                  >
                    <div>
                      <FormCheckbox
                        name='canPayWithCard'
                        formId={formId}
                        size='md'
                        isInvalid={
                          !availablePaymentMethodTypes?.includes('card')
                        }
                      >
                        <div className='text-sm whitespace-nowrap'>
                          Credit or Debit cards
                        </div>
                      </FormCheckbox>
                    </div>
                  </Tooltip>
                  <Tooltip
                    label={
                      availablePaymentMethodTypes?.includes('bacs_debit')
                        ? ''
                        : 'Direct debit not enabled in Stripe'
                    }
                    side='bottom'
                    align='end'
                  >
                    <div>
                      <FormCheckbox
                        name='canPayWithDirectDebit'
                        formId={formId}
                        size='md'
                        isInvalid={
                          !availablePaymentMethodTypes?.includes('ach_debit')
                        }
                      >
                        <div className='text-sm whitespace-nowrap'>
                          Direct Debit via {paymentMethod}
                        </div>
                      </FormCheckbox>
                    </div>
                  </Tooltip>
                </div>
              )}
            </div>

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
                  <div className='text-sm font-normal whitespace-nowrap'>
                    Pay online via Stripe
                  </div>
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
                  <div className='text-sm font-normal whitespace-nowrap'>
                    Bank transfer
                  </div>
                }
              />
            </PaymentDetailsPopover>
            <PaymentDetailsPopover
              withNavigation
              content={
                tenantBillingProfile?.check ? '' : 'Check not enabled yet'
              }
            >
              <FormSwitch
                name='check'
                isInvalid={!tenantBillingProfile?.check}
                formId={formId}
                size='sm'
                labelProps={{
                  margin: 0,
                }}
                label={
                  <div className='text-sm font-normal whitespace-nowrap'>
                    Check
                  </div>
                }
              />
            </PaymentDetailsPopover>
          </div>
        </>
      )}

      <ContractUploader contractId={contractId} />
    </ModalBody>
  );
};
