'use client';
import React, { useMemo, useEffect, PropsWithChildren } from 'react';

import { observer } from 'mobx-react-lite';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { cn } from '@ui/utils/cn';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { countryOptions } from '@shared/util/countryOptions';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { DotSingleEmpty } from '@ui/media/icons/DotSingleEmpty';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Modal, ModalContent, ModalOverlay } from '@ui/overlay/Modal/Modal';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
import { useContractModalStatusContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { useEditContractModalStores } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/EditContractModalStores';

interface SubscriptionServiceModalProps extends PropsWithChildren {
  currency?: string;
  contractId: string;
  showNextInvoice?: boolean;
}

export const ModalWithInvoicePreview = observer(
  ({
    children,
    showNextInvoice,
    contractId,
  }: SubscriptionServiceModalProps) => {
    const { serviceFormStore, invoicePreviewList } =
      useEditContractModalStores();

    const {
      isEditModalOpen,
      onChangeModalMode,
      onEditModalClose,
      addressState,
      detailsState,
    } = useContractModalStateContext();
    const client = getGraphQLClient();

    const { nextInvoice } = useContractModalStatusContext();
    const { data: tenantBillingProfile } =
      useTenantBillingProfilesQuery(client);

    useEffect(() => {
      // run simulation when the edit is opened and there is no next invoice
      if (!nextInvoice && isEditModalOpen) {
        serviceFormStore.runSimulation(invoicePreviewList);
      }
    }, [nextInvoice, isEditModalOpen]);

    const billedTo = {
      addressLine1: addressState.values.addressLine1 ?? '',
      addressLine2: addressState.values.addressLine2 ?? '',
      locality: addressState.values.locality ?? '',
      zip: addressState.values.postalCode ?? '',
      country: addressState?.values?.country?.label ?? '',
      email: addressState.values.billingEmail ?? '',
      name: addressState.values?.organizationLegalName ?? '',
      region: addressState.values?.region ?? '',
    };

    const invoicePreviewStaticData = useMemo(
      () => ({
        status: null,
        invoiceNumber: 'INV-003',
        isBilledToFocused: false,
        note: '',
        currency: detailsState.currency?.value,
        billedTo: {
          addressLine1: addressState.values.addressLine1 ?? '',
          addressLine2: addressState.values.addressLine2 ?? '',
          locality: addressState.values.locality ?? '',
          zip: addressState.values.postalCode ?? '',
          country: addressState?.values?.country?.label ?? '',
          email: addressState.values.billingEmail ?? '',
          name: addressState.values?.organizationLegalName ?? '',
          region: addressState.values?.region ?? '',
        },
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
      [tenantBillingProfile?.tenantBillingProfiles?.[0]],
    );

    return (
      <Modal open={isEditModalOpen} onOpenChange={onEditModalClose}>
        <ModalOverlay />
        <ModalContent
          placement='center'
          className='border-r-2 flex gap-6 bg-transparent shadow-none border-none z-[999] w-full '
          style={{
            minWidth: '1048px',
            minHeight: '80vh',
            boxShadow: 'none',
          }}
        >
          {children}

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
                      canPayWithBankTransfer={true}
                      check={true}
                      availableBankAccount={null}
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
                      onClick={() => invoicePreviewList.setPreviewedInvoice(i)}
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
                              onChangeModalMode(EditModalMode.BillingDetails)
                            }
                            isBilledToFocused={false}
                            shouldBlurDummy={false}
                            note={invoice?.note}
                            invoiceNumber={invoice?.invoiceNumber ?? ''}
                            currency={
                              invoice?.currency ?? detailsState.currency?.value
                            }
                            billedTo={billedTo}
                            from={
                              invoice
                                ? {
                                    addressLine1:
                                      invoice.provider.addressLine1 ?? '',
                                    addressLine2:
                                      invoice.provider.addressLine2 ?? '',
                                    locality:
                                      invoice.provider.addressLocality ?? '',
                                    zip: invoice.provider.addressZip ?? '',
                                    country:
                                      invoice.provider.addressCountry ?? '',
                                    email: '',
                                    name: invoice.provider.name ?? '',
                                    region:
                                      invoice.provider.addressRegion ?? '',
                                  }
                                : invoicePreviewStaticData.from
                            }
                            invoicePeriodStart={invoice?.invoicePeriodStart}
                            invoicePeriodEnd={invoice?.invoicePeriodEnd}
                            tax={invoice.taxDue}
                            lines={invoice?.invoiceLineItems ?? []}
                            subtotal={invoice?.subtotal ?? 10}
                            issueDate={invoice?.issued ?? new Date()}
                            dueDate={invoice?.due ?? new Date()}
                            total={invoice?.total ?? 10}
                            canPayWithBankTransfer={true}
                            check={true}
                            availableBankAccount={null}
                          />
                        </div>
                      </div>
                    );
                  },
                )}
            </div>
          )}
        </ModalContent>
      </Modal>
    );
  },
);
