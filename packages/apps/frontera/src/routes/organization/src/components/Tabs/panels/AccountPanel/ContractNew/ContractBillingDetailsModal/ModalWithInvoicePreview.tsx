import { useMemo, useEffect, PropsWithChildren } from 'react';

import { observer } from 'mobx-react-lite';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { cn } from '@ui/utils/cn';
import { BankAccount } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { countryOptions } from '@shared/util/countryOptions';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { DotSingleEmpty } from '@ui/media/icons/DotEmptySingle.tsx';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { useContractModalStatusContext } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { useEditContractModalStores } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/EditContractModalStores';

interface SubscriptionServiceModalProps extends PropsWithChildren {
  currency?: string;
  showNextInvoice?: boolean;
  allowCheck?: boolean | null;
  billingEnabled?: boolean | null;
  allowAutoPayment?: boolean | null;
  allowBankTransfer?: boolean | null;
  allowOnlinePayment?: boolean | null;
  availableBankAccount?: Partial<BankAccount> | null;
}

export const ModalWithInvoicePreview = observer(
  ({
    children,
    showNextInvoice,
    currency,
    allowCheck,
    allowBankTransfer,
    availableBankAccount,
    billingEnabled,
  }: SubscriptionServiceModalProps) => {
    const { serviceFormStore, invoicePreviewList } =
      useEditContractModalStores();
    const {
      isEditModalOpen,
      onChangeModalMode,
      onEditModalClose,
      addressState,
    } = useContractModalStateContext();
    const client = getGraphQLClient();

    const { nextInvoice } = useContractModalStatusContext();
    const { data: tenantBillingProfile } =
      useTenantBillingProfilesQuery(client);

    useEffect(() => {
      // run simulation when the edit is opened and there is no next invoice

      if (!isEditModalOpen) {
        invoicePreviewList.resetSimulatedInvoices();
      }
    }, [nextInvoice, isEditModalOpen]);

    const billedTo = {
      addressLine1: addressState?.addressLine1 ?? '',
      addressLine2: addressState?.addressLine2 ?? '',
      locality: addressState?.locality ?? '',
      zip: addressState?.postalCode ?? '',
      country: addressState?.country?.label ?? '',
      email: addressState?.billingEmail ?? '',
      name: addressState?.organizationLegalName ?? '',
      region: addressState?.region ?? '',
    };

    const invoicePreviewStaticData = useMemo(
      () => ({
        status: null,
        invoiceNumber: 'INV-003',
        isBilledToFocused: false,
        note: '',
        currency: currency,

        lines: [],
        tax: 0,
        total: 0,
        dueDate: new Date().toISOString(),
        subtotal: 0,
        issueDate: new Date().toISOString(),

        from: tenantBillingProfile?.tenantBillingProfiles?.[0]
          ? {
              addressLine1:
                tenantBillingProfile?.tenantBillingProfiles?.[0]
                  ?.addressLine1 ?? '',
              addressLine2:
                tenantBillingProfile?.tenantBillingProfiles?.[0].addressLine2,
              locality:
                tenantBillingProfile?.tenantBillingProfiles?.[0]?.locality ??
                '',
              zip: tenantBillingProfile?.tenantBillingProfiles?.[0]?.zip ?? '',
              country: tenantBillingProfile?.tenantBillingProfiles?.[0].country
                ? countryOptions.find(
                    (country) =>
                      country.value ===
                      tenantBillingProfile?.tenantBillingProfiles?.[0]?.country,
                  )?.label
                : '',
              email:
                tenantBillingProfile?.tenantBillingProfiles?.[0]
                  ?.sendInvoicesFrom,
              name: tenantBillingProfile?.tenantBillingProfiles?.[0]?.legalName,
              region:
                tenantBillingProfile?.tenantBillingProfiles?.[0]?.region ?? '',
            }
          : {
              addressLine1: '29 Maple Lane',
              addressLine2: 'Springfield, Haven County',
              locality: 'San Francisco',
              zip: '89302',
              country: 'United States of America',
              email: 'invoices@acme.com',
              name: 'Acme Corp.',
              region: 'California',
            },
      }),
      [tenantBillingProfile?.tenantBillingProfiles?.[0], currency],
    );

    return (
      <Modal open={isEditModalOpen} onOpenChange={onEditModalClose}>
        <ModalPortal>
          <ModalOverlay className='z-50' />
          <ModalContent
            placement='center'
            className='border-r-2 flex gap-6 bg-transparent shadow-none border-none z-[999] w-full '
            style={{
              minWidth: billingEnabled ? '1048px' : 'auto',
              minHeight: '80vh',
              boxShadow: 'none',
            }}
          >
            {children}

            {billingEnabled && (
              <>
                {showNextInvoice === undefined && (
                  <div
                    style={{ minWidth: '600px' }}
                    className={cn('bg-white rounded  h-full absolute z-10')}
                  >
                    <InvoiceSkeleton />
                  </div>
                )}
                {showNextInvoice && (
                  <div className='h-auto w-full flex relative min-w-[600px]'>
                    {invoicePreviewList.simulatedInvoices.length === 0 && (
                      <div
                        style={{ minWidth: '600px' }}
                        className={cn('bg-white rounded  h-full absolute z-1 ')}
                      >
                        <div className='w-full h-full'>
                          <Invoice
                            {...invoicePreviewStaticData}
                            shouldBlurDummy={true}
                            onOpenAddressDetailsModal={() =>
                              onChangeModalMode(EditModalMode.BillingDetails)
                            }
                            billedTo={billedTo}
                            currency={currency}
                            canPayWithBankTransfer={allowBankTransfer}
                            check={allowCheck}
                            availableBankAccount={availableBankAccount}
                          />
                        </div>
                      </div>
                    )}

                    {invoicePreviewList.simulatedInvoices.length > 1 && (
                      <div className='absolute top-[-30px] right-0 text-white text-base '>
                        <IconButton
                          variant='ghost'
                          className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'
                          aria-label='prev'
                          icon={<ChevronLeft className='text-inherit' />}
                          onClick={() => {
                            if (invoicePreviewList.previewedInvoiceIndex <= 0) {
                              return;
                            }
                            invoicePreviewList.setPreviewedInvoice(
                              invoicePreviewList.previewedInvoiceIndex - 1,
                            );
                          }}
                        />
                        {invoicePreviewList.simulatedInvoices.map((e, i) => (
                          <IconButton
                            variant='ghost'
                            className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'
                            key={e.invoice?.invoiceNumber}
                            aria-label='prev'
                            onClick={() =>
                              invoicePreviewList.setPreviewedInvoice(i)
                            }
                            icon={
                              i === invoicePreviewList.previewedInvoiceIndex ? (
                                <DotSingle className='text-inherit' />
                              ) : (
                                <DotSingleEmpty className='text-inherit' />
                              )
                            }
                          />
                        ))}

                        <IconButton
                          variant='ghost'
                          className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'
                          aria-label='prev'
                          icon={<ChevronRight className='text-inherit' />}
                          onClick={() => {
                            if (
                              invoicePreviewList.previewedInvoiceIndex >=
                              invoicePreviewList.simulatedInvoices.length - 1
                            ) {
                              return;
                            }
                            invoicePreviewList.setPreviewedInvoice(
                              invoicePreviewList.previewedInvoiceIndex + 1,
                            );
                          }}
                        />
                      </div>
                    )}

                    {serviceFormStore.isSimulationRunning && (
                      <div
                        style={{ minWidth: '600px' }}
                        className={cn('bg-white rounded  h-full absolute z-10')}
                      >
                        <InvoiceSkeleton />
                      </div>
                    )}

                    {!serviceFormStore.isSimulationRunning &&
                      invoicePreviewList.simulatedInvoices.map(
                        ({ invoice }, index) => {
                          if (!invoice) return null;

                          return (
                            <div
                              style={{ minWidth: '600px' }}
                              className={cn(
                                'bg-white rounded  h-full absolute z-1 ',
                                {
                                  '-rotate-[1.15deg] shadow-lg  ':
                                    index !==
                                    invoicePreviewList.previewedInvoiceIndex,
                                  'shadow-md z-10 animate-zIndex':
                                    index ===
                                    invoicePreviewList.previewedInvoiceIndex,
                                },
                              )}
                              key={`${invoice?.invoiceNumber}-${invoice.total}-invoice-preview`}
                            >
                              <div className='w-full h-full'>
                                <Invoice
                                  onOpenAddressDetailsModal={() =>
                                    onChangeModalMode(
                                      EditModalMode.BillingDetails,
                                    )
                                  }
                                  isBilledToFocused={false}
                                  shouldBlurDummy={false}
                                  note={invoice?.note}
                                  invoiceNumber={invoice?.invoiceNumber ?? ''}
                                  currency={invoice?.currency ?? currency}
                                  billedTo={billedTo}
                                  from={
                                    invoice
                                      ? {
                                          addressLine1:
                                            invoice.provider.addressLine1 ?? '',
                                          addressLine2:
                                            invoice.provider.addressLine2 ?? '',
                                          locality:
                                            invoice.provider.addressLocality ??
                                            '',
                                          zip:
                                            invoice.provider.addressZip ?? '',
                                          country:
                                            invoice.provider.addressCountry ??
                                            '',
                                          email: '',
                                          name: invoice.provider.name ?? '',
                                          region:
                                            invoice.provider.addressRegion ??
                                            '',
                                        }
                                      : invoicePreviewStaticData.from
                                  }
                                  invoicePeriodStart={
                                    invoice?.invoicePeriodStart
                                  }
                                  invoicePeriodEnd={invoice?.invoicePeriodEnd}
                                  tax={invoice.taxDue}
                                  lines={invoice?.invoiceLineItems ?? []}
                                  subtotal={invoice?.subtotal ?? 10}
                                  issueDate={invoice?.issued ?? new Date()}
                                  dueDate={invoice?.due ?? new Date()}
                                  total={invoice?.total ?? 10}
                                  canPayWithBankTransfer={true}
                                  check={allowCheck}
                                  availableBankAccount={availableBankAccount}
                                />
                              </div>
                            </div>
                          );
                        },
                      )}
                  </div>
                )}
              </>
            )}
          </ModalContent>
        </ModalPortal>
      </Modal>
    );
  },
);
