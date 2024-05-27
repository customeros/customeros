import { FC, useMemo } from 'react';
import { useField } from 'react-inverted-form';

import { toZonedTime } from 'date-fns-tz';
import { useConnections } from '@integration-app/react';
import { useGetExternalSystemInstancesQuery } from '@settings/graphql/getExternalSystemInstances.generated';

import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { FormSwitch } from '@ui/form/Switch/FormSwitch';
import { Divider } from '@ui/presentation/Divider/Divider';
import { FormCheckbox } from '@ui/form/Checkbox/FormCheckbox';
import { currencyOptions } from '@shared/util/currencyOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import {
  Currency,
  BankAccount,
  ContractStatus,
  ExternalSystemType,
  TenantBillingProfile,
} from '@graphql/types';
import {
  paymentDueOptions,
  contractBillingCycleOptions,
} from '@organization/components/Tabs/panels/AccountPanel/utils';

import { Services } from './Services';
import { InlineSelect } from './InlineSelect';
import { ContractUploader } from './ContractUploader';
import { CommittedPeriodInput } from './CommittedPeriodInput';
import { PaymentDetailsPopover } from './PaymentDetailsPopover';

interface SubscriptionServiceModalProps {
  formId: string;
  currency?: string;
  renewedAt?: string;
  contractId: string;
  billingEnabled?: boolean;
  payAutomatically?: boolean | null;
  contractStatus?: ContractStatus | null;
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
  billingEnabled,
  renewedAt,
  contractStatus,
}) => {
  const client = getGraphQLClient();
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
  const { getInputProps: billingEnabledInputProps } = useField(
    'billingEnabled',
    formId,
  );
  const { onChange: onChangeBillingEnabled, value: billingEnabledValue } =
    billingEnabledInputProps();
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
  const renewalDate = renewedAt
    ? DateTimeUtils.format(
        toZonedTime(renewedAt, 'UTC').toUTCString(),
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;

  return (
    <ModalBody className='flex flex-col flex-1 p-0'>
      <ul className='mb-2 list-disc ml-5'>
        <li className='text-base '>
          <div className='flex items-baseline'>
            <CommittedPeriodInput formId={formId} />

            <span className='whitespace-nowrap mr-1'>contract, starting </span>

            <DatePickerUnderline formId={formId} name='serviceStarted' />
          </div>
        </li>
        <li className='text-base mt-1.5'>
          <div className='flex items-baseline'>
            Live until {renewalDate},{' '}
            <Button
              variant='ghost'
              size='sm'
              className='font-normal text-base p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline'
              onClick={() => onChangeAutoRenew(!autoRenewValue)}
            >
              {autoRenewValue ? 'auto-renews' : 'not auto-renewing'}
            </Button>
          </div>
        </li>
        <li className='text-base '>
          <div className='flex items-baseline'>
            <span className='whitespace-nowrap'>Contracting in</span>
            <div>
              <InlineSelect
                label='Currency'
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
      <Services currency={currency} contractStatus={contractStatus} />
      {billingEnabled && (
        <>
          <div className='flex relative items-center h-8 mb-1'>
            <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
              Billing policy
            </p>
            <Divider />
          </div>
          <ul className='mb-2 list-disc ml-5'>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap mr-1'>Billing starts </span>

                <DatePickerUnderline formId={formId} name='invoicingStarted' />
              </div>
            </li>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap'>Invoices are sent</span>
                <div>
                  <InlineSelect
                    label='billing period'
                    placeholder='billing period'
                    name='billingCycle'
                    formId={formId}
                    options={contractBillingCycleOptions}
                    size='xs'
                  />
                </div>
                <span className='whitespace-nowrap ml-0.5'>
                  on the billing start day
                </span>
              </div>
            </li>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap '>Customer has</span>
                <div>
                  <InlineSelect
                    label='Payment due'
                    placeholder='0 days'
                    name='dueDays'
                    formId={formId}
                    options={paymentDueOptions}
                    size='xs'
                  />
                </div>
                <span className='whitespace-nowrap ml-0.5'>to pay</span>
              </div>
            </li>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap '>Billing is </span>
                <div>
                  <Button
                    variant='ghost'
                    size='sm'
                    className='font-normal text-base p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline'
                    onClick={() => onChangeBillingEnabled(!billingEnabledValue)}
                  >
                    {billingEnabledValue ? 'enabled' : 'disabled'}
                  </Button>
                </div>
                <span className='whitespace-nowrap ml-0.5'>in CustomerOS</span>
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
              <FormSwitch
                name='payAutomatically'
                formId={formId}
                isInvalid={!isStripeActive}
                size='sm'
                labelProps={{
                  className: 'm-0',
                }}
                label={
                  <PaymentDetailsPopover
                    content={
                      isStripeActive ? '' : 'No payment provider enabled'
                    }
                    withNavigation
                  >
                    <div className='text-base font-normal whitespace-nowrap'>
                      Auto-payment via Stripe
                    </div>
                  </PaymentDetailsPopover>
                }
              />
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
                        size='sm'
                        iconSize='sm'
                        disabled={
                          !availablePaymentMethodTypes?.includes('card')
                        }
                      >
                        <div className='text-base whitespace-nowrap'>
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
                        size='sm'
                        iconColorScheme='warning'
                        iconSize='sm'
                        colorScheme='warning'
                        disabled={
                          !availablePaymentMethodTypes?.includes('ach_debit')
                        }
                      >
                        <div className='text-base whitespace-nowrap'>
                          Direct Debit via {paymentMethod}
                        </div>
                      </FormCheckbox>
                    </div>
                  </Tooltip>
                </div>
              )}
            </div>
            <FormSwitch
              name='payOnline'
              formId={formId}
              isInvalid={!isStripeActive}
              size='sm'
              labelProps={{
                className: 'm-0',
              }}
              label={
                <PaymentDetailsPopover
                  content={isStripeActive ? '' : 'No payment provider enabled'}
                  withNavigation
                >
                  <div className='text-base font-normal whitespace-nowrap'>
                    Pay online via Stripe
                  </div>
                </PaymentDetailsPopover>
              }
            />
            <FormSwitch
              name='canPayWithBankTransfer'
              isInvalid={!!bankTransferPopoverContent.length}
              formId={formId}
              size='sm'
              labelProps={{
                className: 'm-0',
              }}
              label={
                <PaymentDetailsPopover
                  withNavigation
                  content={bankTransferPopoverContent}
                >
                  <div className='text-base font-normal whitespace-nowrap'>
                    Bank transfer
                  </div>
                </PaymentDetailsPopover>
              }
            />

            <FormSwitch
              name='check'
              isInvalid={!tenantBillingProfile?.check}
              formId={formId}
              size='sm'
              labelProps={{
                className: 'm-0',
              }}
              label={
                <PaymentDetailsPopover
                  withNavigation
                  content={
                    tenantBillingProfile?.check ? '' : 'Check not enabled yet'
                  }
                >
                  <div className='text-base font-normal whitespace-nowrap'>
                    Check
                  </div>
                </PaymentDetailsPopover>
              }
            />
          </div>
        </>
      )}

      <ContractUploader contractId={contractId} />
    </ModalBody>
  );
};
