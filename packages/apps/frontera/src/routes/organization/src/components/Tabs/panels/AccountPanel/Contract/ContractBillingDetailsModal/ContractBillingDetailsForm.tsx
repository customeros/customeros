import { FC, useMemo } from 'react';

import { Store } from '@store/store.ts';
import { toZonedTime } from 'date-fns-tz';
import { useConnections } from '@integration-app/react';

import { Switch } from '@ui/form/Switch';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.tsx';
import { Divider } from '@ui/presentation/Divider/Divider';
import { currencyOptions } from '@shared/util/currencyOptions';
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
  renewedAt?: string;
  contractId: string;
  billingEnabled?: boolean;
  payAutomatically?: boolean | null;
  contractStatus?: ContractStatus | null;
  tenantBillingProfile?: TenantBillingProfile | null;
  bankAccounts: Array<Store<BankAccount>> | null | undefined;
}

export const ContractBillingDetailsForm: FC<SubscriptionServiceModalProps> = ({
  contractId,
  tenantBillingProfile,
  bankAccounts,
  payAutomatically,
  billingEnabled,
  renewedAt,
  contractStatus,
}) => {
  const store = useStore();
  const externalSystemInstances = store.externalSystemInstances;
  const contractStore = store.contracts.value.get(contractId);
  const currency = contractStore?.value?.currency;
  const availablePaymentMethodTypes = externalSystemInstances?.value?.find(
    (e) => e.type === ExternalSystemType.Stripe,
  )?.stripeDetails?.paymentMethodTypes;
  const { items: iConnections } = useConnections();
  const isStripeActive = !!iConnections
    .map((item) => item.integration?.key)
    .find((e) => e === 'stripe');

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
      (account) => account?.value?.currency === currency,
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
            <CommittedPeriodInput contractId={contractId} />

            <span className='whitespace-nowrap mr-1'>contract, starting </span>

            <DatePickerUnderline
              value={contractStore?.value?.serviceStarted}
              onChange={(date) =>
                contractStore?.update(
                  (prev) => ({
                    ...prev,
                    serviceStarted: date,
                  }),
                  { mutate: false },
                )
              }
            />
          </div>
        </li>
        <li className='text-base mt-1.5'>
          <div className='flex items-baseline'>
            Live until {renewalDate},{' '}
            <Button
              variant='ghost'
              size='sm'
              className='font-normal text-base p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline'
              onClick={() =>
                contractStore?.update(
                  (prev) => ({
                    ...prev,
                    autoRenew: !contractStore?.value.autoRenew,
                  }),
                  { mutate: false },
                )
              }
            >
              {contractStore?.value.autoRenew
                ? 'auto-renews'
                : 'not auto-renewing'}
            </Button>
          </div>
        </li>
        <li className='text-base '>
          <div className='flex items-baseline'>
            <span className='whitespace-nowrap'>Contracting in</span>
            <div className='z-30'>
              <InlineSelect
                id='contract-currency'
                name='contract-currency'
                label='Currency'
                placeholder='Invoice currency'
                value={currency}
                onChange={(selectedOption) =>
                  contractStore?.update(
                    (contract) => ({
                      ...contract,
                      currency: selectedOption.value as Currency,
                    }),
                    { mutate: false },
                  )
                }
                options={currencyOptions}
                size='xs'
              />
            </div>
          </div>
        </li>
      </ul>
      <Services
        id={contractId}
        currency={currency}
        contractStatus={contractStatus}
        billingEnabled={billingEnabled}
      />
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

                <DatePickerUnderline
                  value={contractStore?.value?.billingDetails?.invoicingStarted}
                  onChange={(date) =>
                    contractStore?.update(
                      (prev) => ({
                        ...prev,
                        billingDetails: {
                          ...prev.billingDetails,
                          invoicingStarted: date,
                        },
                      }),
                      { mutate: false },
                    )
                  }
                />
              </div>
            </li>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap'>Invoices are sent</span>
                <span className='z-20'>
                  <InlineSelect
                    label='billing period'
                    placeholder='billing period'
                    id='contract-billingCycle'
                    name='contract-billingCycle'
                    options={contractBillingCycleOptions}
                    value={
                      contractStore?.value?.billingDetails?.billingCycleInMonths
                    }
                    onChange={(selectedOption) =>
                      contractStore?.update(
                        (contract) => ({
                          ...contract,
                          billingDetails: {
                            ...contract.billingDetails,
                            billingCycleInMonths: selectedOption.value,
                          },
                        }),
                        { mutate: false },
                      )
                    }
                    size='xs'
                  />
                </span>
                <span className='whitespace-nowrap ml-0.5'>
                  on the billing start day
                </span>
              </div>
            </li>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <span className='whitespace-nowrap '>Customer has</span>
                <div className='inline z-10'>
                  <InlineSelect
                    label='Payment due'
                    placeholder='0 days'
                    name='dueDays'
                    id='dueDays'
                    value={contractStore?.value?.billingDetails?.dueDays}
                    onChange={(selectedOption) =>
                      contractStore?.update(
                        (contract) => ({
                          ...contract,
                          billingDetails: {
                            ...contract.billingDetails,
                            dueDays: selectedOption.value,
                          },
                        }),
                        { mutate: false },
                      )
                    }
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
                    onClick={() =>
                      contractStore?.update(
                        (contract) => ({
                          ...contract,
                          billingEnabled: !contractStore?.value.billingEnabled,
                        }),
                        { mutate: false },
                      )
                    }
                  >
                    {contractStore?.value.billingEnabled
                      ? 'enabled'
                      : 'disabled'}
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
            <div className='flex w-full justify-between items-center'>
              <PaymentDetailsPopover
                content={isStripeActive ? '' : 'No payment provider enabled'}
                withNavigation
              >
                <div className='text-base font-normal whitespace-nowrap'>
                  Pay online via Stripe
                </div>
              </PaymentDetailsPopover>
              <Switch
                name='payAutomatically'
                isInvalid={!isStripeActive}
                size='sm'
                isChecked={
                  !!contractStore?.value?.billingDetails?.payAutomatically
                }
                onChange={(value) =>
                  contractStore?.update(
                    (contract) => ({
                      ...contract,
                      billingDetails: {
                        ...contract.billingDetails,
                        payAutomatically: value,
                      },
                    }),
                    { mutate: false },
                  )
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
                      <Checkbox
                        key='canPayWithCard'
                        disabled={
                          !availablePaymentMethodTypes?.includes('card')
                        }
                        isChecked={
                          !!contractStore?.value.billingDetails?.canPayWithCard
                        }
                        onChange={(value) =>
                          contractStore?.update(
                            (contract) => ({
                              ...contract,
                              billingDetails: {
                                ...contract.billingDetails,
                                canPayWithCard: !!value,
                              },
                            }),
                            { mutate: false },
                          )
                        }
                      >
                        <div className='text-base whitespace-nowrap'>
                          Credit or Debit cards
                        </div>
                      </Checkbox>
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
                      <Checkbox
                        key='canPayWithDirectDebit'
                        disabled={
                          !availablePaymentMethodTypes?.includes('ach_debit')
                        }
                        isChecked={
                          !!contractStore?.value.billingDetails
                            ?.canPayWithDirectDebit
                        }
                        onChange={(value) =>
                          contractStore?.update(
                            (contract) => ({
                              ...contract,
                              billingDetails: {
                                ...contract.billingDetails,
                                canPayWithDirectDebit: !!value,
                              },
                            }),
                            { mutate: false },
                          )
                        }
                      >
                        <div className='text-base whitespace-nowrap'>
                          Direct Debit via {paymentMethod}
                        </div>
                      </Checkbox>
                    </div>
                  </Tooltip>
                </div>
              )}
            </div>

            <div className='flex w-full justify-between items-center'>
              <PaymentDetailsPopover
                content={isStripeActive ? '' : 'No payment provider enabled'}
                withNavigation
              >
                <div className='text-base font-normal whitespace-nowrap'>
                  Pay online via Stripe
                </div>
              </PaymentDetailsPopover>
              <Switch
                name='payOnline'
                isInvalid={!!bankTransferPopoverContent.length}
                size='sm'
                isChecked={!!contractStore?.value?.billingDetails?.payOnline}
                onChange={(value) =>
                  contractStore?.update(
                    (contract) => ({
                      ...contract,
                      billingDetails: {
                        ...contract.billingDetails,
                        payOnline: value,
                      },
                    }),
                    { mutate: false },
                  )
                }
              />
            </div>

            <div className='flex w-full justify-between items-center'>
              <PaymentDetailsPopover
                withNavigation
                content={bankTransferPopoverContent}
              >
                <div className='text-base font-normal whitespace-nowrap'>
                  Bank transfer
                </div>
              </PaymentDetailsPopover>
              <Switch
                name='canPayWithBankTransfer'
                isInvalid={!!bankTransferPopoverContent.length}
                size='sm'
                isChecked={
                  !!contractStore?.value?.billingDetails?.canPayWithBankTransfer
                }
                onChange={(value) =>
                  contractStore?.update(
                    (contract) => ({
                      ...contract,
                      billingDetails: {
                        ...contract.billingDetails,
                        canPayWithBankTransfer: value,
                      },
                    }),
                    { mutate: false },
                  )
                }
              />
            </div>

            <div className='flex w-full justify-between items-center'>
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
              <Switch
                name='check'
                isInvalid={!tenantBillingProfile?.check}
                size='sm'
                isChecked={!!contractStore?.value?.billingDetails?.check}
                onChange={(value) =>
                  contractStore?.update(
                    (contract) => ({
                      ...contract,
                      billingDetails: {
                        ...contract.billingDetails,
                        check: value,
                      },
                    }),
                    { mutate: false },
                  )
                }
              />
            </div>
          </div>
        </>
      )}

      <ContractUploader contractId={contractId} />
    </ModalBody>
  );
};
