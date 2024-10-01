import { useMemo } from 'react';

import { Store } from '@store/store.ts';
import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';
import { useConnections } from '@integration-app/react';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { Switch } from '@ui/form/Switch';
import { DateTimeUtils } from '@utils/date.ts';
import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { Button } from '@ui/form/Button/Button.tsx';
import { ModalBody } from '@ui/overlay/Modal/Modal.tsx';
import { Divider } from '@ui/presentation/Divider/Divider.tsx';
import { currencyOptions } from '@shared/util/currencyOptions.ts';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline.tsx';
import {
  Currency,
  BankAccount,
  ContractStatus,
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

export const ContractBillingDetailsForm = observer(
  ({
    contractId,
    tenantBillingProfile,
    bankAccounts,
    billingEnabled,
    contractStatus,
    openAddressModal,
  }: SubscriptionServiceModalProps) => {
    const store = useStore();
    const contractStore = store.contracts.value.get(
      contractId,
    ) as ContractStore;

    const currency = contractStore?.tempValue?.currency;

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
                      contractStore?.tempValue?.billingDetails?.invoicingStarted
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
                      value={contractStore?.tempValue?.billingDetails?.dueDays}
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
                      dataTest='contract-billing-details-address'
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

            <div className='flex flex-col gap-2 mb-2'>
              <div className='flex flex-col w-full justify-between items-start'>
                <div className='flex  items-center justify-between w-full'>
                  <PaymentDetailsPopover
                    withNavigation
                    content={
                      isStripeActive ? '' : 'No payment provider enabled'
                    }
                  >
                    <div className='text-sm font-normal whitespace-nowrap'>
                      Auto-charge card via Stripe
                    </div>
                  </PaymentDetailsPopover>

                  <Switch
                    size='sm'
                    name='payAutomatically'
                    isInvalid={!isStripeActive}
                    isChecked={
                      !!contractStore?.tempValue?.billingDetails?.payOnline
                    }
                    onChange={(value) => {
                      contractStore?.updateTemp((contract) => ({
                        ...contract,
                        billingDetails: {
                          ...contract.billingDetails,
                          payOnline: value,
                          payAutomatically: value,
                        },
                      }));
                    }}
                  />
                </div>
              </div>

              {contractStore?.tempValue.billingDetails?.payOnline && (
                <RadioGroup
                  name='created-date'
                  value={`${!!contractStore.tempValue.billingDetails
                    ?.payAutomatically}`}
                  onValueChange={(newValue) => {
                    contractStore?.updateTemp((contract) => ({
                      ...contract,
                      billingDetails: {
                        ...contract.billingDetails,
                        payAutomatically: newValue === 'true',
                      },
                    }));
                  }}
                >
                  <div className='flex flex-col gap-2 items-start'>
                    <Radio value={'true'}>
                      <span>Auto-charge card</span>
                    </Radio>
                    <Radio value={'false'}>
                      <span>One-off payment link</span>
                    </Radio>
                  </div>
                </RadioGroup>
              )}

              <div className='flex w-full justify-between items-center'>
                <PaymentDetailsPopover
                  withNavigation
                  content={bankTransferPopoverContent}
                >
                  <div className='font-normal whitespace-nowrap'>
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
                  <div className='font-normal whitespace-nowrap'>Checks</div>
                </PaymentDetailsPopover>
                <Switch
                  size='sm'
                  name='check'
                  isInvalid={!tenantBillingProfile?.check}
                  isChecked={!!contractStore?.tempValue?.billingDetails?.check}
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
