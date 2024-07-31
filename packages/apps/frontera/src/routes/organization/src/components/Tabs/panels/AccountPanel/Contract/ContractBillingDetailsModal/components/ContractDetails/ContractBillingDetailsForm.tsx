import { FC, useMemo } from 'react';

import { Store } from '@store/store.ts';
import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';
import { useConnections } from '@integration-app/react';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { Switch } from '@ui/form/Switch';
import { DateTimeUtils } from '@utils/date.ts';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { ModalBody } from '@ui/overlay/Modal/Modal.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.tsx';
import { Divider } from '@ui/presentation/Divider/Divider.tsx';
import { currencyOptions } from '@shared/util/currencyOptions.ts';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline.tsx';
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
} from '@organization/components/Tabs/panels/AccountPanel/utils.ts';

import { Services } from '../Services';
import { InlineSelect } from './InlineSelect.tsx';
import { ContractUploader } from './ContractUploader.tsx';
import { CommittedPeriodInput } from './CommittedPeriodInput.tsx';
import { PaymentDetailsPopover } from './PaymentDetailsPopover.tsx';

interface SubscriptionServiceModalProps {
  contractId: string;
  billingEnabled?: boolean;
  openAddressModal: () => void;
  contractStatus?: ContractStatus | null;
  tenantBillingProfile?: TenantBillingProfile | null;
  bankAccounts: Array<Store<BankAccount>> | null | undefined;
}

export const ContractBillingDetailsForm: FC<SubscriptionServiceModalProps> =
  observer(
    ({
      contractId,
      tenantBillingProfile,
      bankAccounts,
      billingEnabled,
      contractStatus,
      openAddressModal,
    }) => {
      const store = useStore();
      const externalSystemInstances = store.externalSystemInstances;
      const contractStore = store.contracts.value.get(
        contractId,
      ) as ContractStore;

      const currency = contractStore?.tempValue?.currency;
      const availablePaymentMethodTypes = externalSystemInstances?.value?.find(
        (e) => e.type === ExternalSystemType.Stripe,
      )?.stripeDetails?.paymentMethodTypes;
      const { items: iConnections } = useConnections();
      const isStripeActive = !!iConnections
        .map((item) => item.integration?.key)
        .find((e) => e === 'stripe');

      const payAutomatically =
        contractStore?.tempValue?.billingDetails?.payAutomatically;

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

      const renewalCalculatedDate = useMemo(() => {
        if (!contractStore?.tempValue?.serviceStarted) return null;
        const parsed = contractStore?.tempValue?.committedPeriodInMonths
          ? parseFloat(contractStore?.tempValue?.committedPeriodInMonths)
          : 1;

        return DateTimeUtils.addMonth(
          contractStore.tempValue.serviceStarted,
          parsed,
        );
      }, [
        contractStore?.tempValue?.serviceStarted,
        contractStore?.tempValue?.committedPeriodInMonths,
      ]);

      return (
        <ModalBody className='flex flex-col flex-1 p-0'>
          <ul className='mb-2 list-disc ml-5'>
            <li className='text-base '>
              <div className='flex items-baseline'>
                <CommittedPeriodInput contractId={contractId} />

                <span className='whitespace-nowrap mr-1'>
                  contract, starting{' '}
                </span>

                <DatePickerUnderline
                  value={toZonedTime(
                    contractStore?.tempValue?.serviceStarted,
                    'UTC',
                  )}
                  onChange={(date) =>
                    contractStore?.updateTemp((prev) => ({
                      ...prev,
                      serviceStarted: date,
                    }))
                  }
                />
              </div>
            </li>
            <li className='text-base'>
              <div className='flex items-baseline'>
                Live until{' '}
                {renewalCalculatedDate
                  ? DateTimeUtils.format(
                      toZonedTime(renewalCalculatedDate, 'UTC').toUTCString(),
                      DateTimeUtils.dateWithAbreviatedMonth,
                    )
                  : '...'}
                ,{' '}
                <Button
                  size='sm'
                  variant='ghost'
                  className='font-normal text-base p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline'
                  onClick={() =>
                    contractStore?.updateTemp((prev) => ({
                      ...prev,
                      autoRenew: !contractStore?.tempValue.autoRenew,
                    }))
                  }
                >
                  {contractStore?.tempValue.autoRenew
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
                    size='xs'
                    label='Currency'
                    value={currency}
                    id='contract-currency'
                    name='contract-currency'
                    options={currencyOptions}
                    placeholder='Invoice currency'
                    onChange={(selectedOption) =>
                      contractStore?.updateTemp((contract) => ({
                        ...contract,
                        currency: selectedOption.value as Currency,
                      }))
                    }
                  />
                </div>
              </div>
            </li>
          </ul>
          <Services
            id={contractId}
            contractStatus={contractStatus}
            currency={currency ?? Currency.Usd}
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
                    <span className='whitespace-nowrap mr-1'>
                      Billing starts{' '}
                    </span>

                    <DatePickerUnderline
                      value={
                        contractStore?.tempValue?.billingDetails
                          ?.invoicingStarted
                      }
                      onChange={(date) =>
                        contractStore?.updateTemp((prev) => ({
                          ...prev,
                          billingDetails: {
                            ...prev.billingDetails,
                            invoicingStarted: date,
                          },
                        }))
                      }
                    />
                  </div>
                </li>
                <li className='text-base '>
                  <div className='flex items-baseline'>
                    <span className='whitespace-nowrap'>Invoices are sent</span>
                    <span className='z-20'>
                      <InlineSelect
                        size='xs'
                        label='billing period'
                        id='contract-billingCycle'
                        placeholder='billing period'
                        name='contract-billingCycle'
                        options={contractBillingCycleOptions}
                        value={
                          contractStore?.tempValue?.billingDetails
                            ?.billingCycleInMonths
                        }
                        onChange={(selectedOption) =>
                          contractStore?.updateTemp((contract) => ({
                            ...contract,
                            billingDetails: {
                              ...contract.billingDetails,
                              billingCycleInMonths: selectedOption.value,
                            },
                          }))
                        }
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
                        size='xs'
                        id='dueDays'
                        name='dueDays'
                        label='Payment due'
                        placeholder='0 days'
                        options={paymentDueOptions}
                        value={
                          contractStore?.tempValue?.billingDetails?.dueDays
                        }
                        onChange={(selectedOption) =>
                          contractStore?.updateTemp((contract) => ({
                            ...contract,
                            billingDetails: {
                              ...contract.billingDetails,
                              dueDays: selectedOption.value,
                            },
                          }))
                        }
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
                        size='sm'
                        variant='ghost'
                        className='font-normal text-base p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline'
                        onClick={() =>
                          contractStore?.updateTemp((contract) => ({
                            ...contract,
                            billingEnabled:
                              !contractStore?.tempValue.billingEnabled,
                          }))
                        }
                      >
                        {contractStore?.tempValue.billingEnabled
                          ? 'enabled'
                          : 'disabled'}
                      </Button>
                    </div>
                    <span className='whitespace-nowrap ml-0.5'>
                      in CustomerOS
                    </span>
                  </div>
                </li>
                <li className='text-base '>
                  <div className='flex items-baseline'>
                    <span className='whitespace-nowrap '>
                      Invoices are billed to{' '}
                    </span>
                    <div>
                      <Button
                        size='sm'
                        variant='ghost'
                        onClick={openAddressModal}
                        className='font-normal text-base p-0 ml-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline'
                      >
                        {contractStore?.tempValue.billingDetails
                          ?.organizationLegalName || 'this address'}
                      </Button>
                    </div>
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
                <div className='flex flex-col w-full justify-between items-start'>
                  <div className='flex  items-center justify-between w-full'>
                    <PaymentDetailsPopover
                      withNavigation
                      content={
                        isStripeActive ? '' : 'No payment provider enabled'
                      }
                    >
                      <div className='text-base font-normal whitespace-nowrap'>
                        Auto-payment via Stripe
                      </div>
                    </PaymentDetailsPopover>

                    <Switch
                      size='sm'
                      name='payAutomatically'
                      isInvalid={!isStripeActive}
                      isChecked={
                        !!contractStore?.tempValue?.billingDetails
                          ?.payAutomatically
                      }
                      onChange={(value) => {
                        contractStore?.updateTemp((contract) => ({
                          ...contract,
                          billingDetails: {
                            ...contract.billingDetails,
                            payAutomatically: !!value,
                          },
                        }));
                      }}
                    />
                  </div>

                  {isStripeActive && payAutomatically && (
                    <div className='flex flex-col gap-1 ml-2 mt-1'>
                      <Tooltip
                        align='end'
                        side='bottom'
                        label={
                          availablePaymentMethodTypes?.includes('card')
                            ? ''
                            : 'Credit or Debit card not enabled in Stripe'
                        }
                      >
                        <div>
                          <Checkbox
                            size='sm'
                            key='canPayWithCard'
                            disabled={
                              !availablePaymentMethodTypes?.includes('card')
                            }
                            isChecked={
                              !!contractStore?.tempValue.billingDetails
                                ?.canPayWithCard
                            }
                            onChange={(value) =>
                              contractStore?.updateTemp((contract) => ({
                                ...contract,
                                billingDetails: {
                                  ...contract.billingDetails,
                                  canPayWithCard: !!value,
                                },
                              }))
                            }
                          >
                            <div className='text-base whitespace-nowrap'>
                              Credit or Debit cards
                            </div>
                          </Checkbox>
                        </div>
                      </Tooltip>
                      <Tooltip
                        align='end'
                        side='bottom'
                        label={
                          availablePaymentMethodTypes?.includes('bacs_debit')
                            ? ''
                            : 'Direct debit not enabled in Stripe'
                        }
                      >
                        <div>
                          <Checkbox
                            size='sm'
                            key='canPayWithDirectDebit'
                            isChecked={
                              !!contractStore?.tempValue.billingDetails
                                ?.canPayWithDirectDebit
                            }
                            disabled={
                              !availablePaymentMethodTypes?.includes(
                                'ach_debit',
                              )
                            }
                            onChange={(value) =>
                              contractStore?.updateTemp((contract) => ({
                                ...contract,
                                billingDetails: {
                                  ...contract.billingDetails,
                                  canPayWithDirectDebit: !!value,
                                },
                              }))
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
                    withNavigation
                    content={
                      isStripeActive ? '' : 'No payment provider enabled'
                    }
                  >
                    <div className='text-base font-normal whitespace-nowrap'>
                      Pay online via Stripe
                    </div>
                  </PaymentDetailsPopover>
                  <Switch
                    size='sm'
                    name='payOnline'
                    isInvalid={!!bankTransferPopoverContent.length}
                    isChecked={
                      !!contractStore?.tempValue?.billingDetails?.payOnline
                    }
                    onChange={(value) =>
                      contractStore?.updateTemp((contract) => ({
                        ...contract,
                        billingDetails: {
                          ...contract.billingDetails,
                          payOnline: value,
                        },
                      }))
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
                    size='sm'
                    name='canPayWithBankTransfer'
                    isInvalid={!!bankTransferPopoverContent.length}
                    isChecked={
                      !!contractStore?.tempValue?.billingDetails
                        ?.canPayWithBankTransfer
                    }
                    onChange={(value) =>
                      contractStore?.updateTemp((contract) => ({
                        ...contract,
                        billingDetails: {
                          ...contract.billingDetails,
                          canPayWithBankTransfer: value,
                        },
                      }))
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
                    size='sm'
                    name='check'
                    isInvalid={!tenantBillingProfile?.check}
                    isChecked={
                      !!contractStore?.tempValue?.billingDetails?.check
                    }
                    onChange={(value) =>
                      contractStore?.updateTemp((contract) => ({
                        ...contract,
                        billingDetails: {
                          ...contract.billingDetails,
                          check: value,
                        },
                      }))
                    }
                  />
                </div>
              </div>
            </>
          )}

          <ContractUploader contractId={contractId} />
        </ModalBody>
      );
    },
  );
