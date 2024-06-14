import { useMemo, useEffect, PropsWithChildren } from 'react';

import { toJS } from 'mobx';
import { Store } from '@store/store.ts';
import { observer } from 'mobx-react-lite';
import { isAfter } from 'date-fns/isAfter';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { BankAccount } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Invoice as TInvoice } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { countryOptions } from '@shared/util/countryOptions';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { DotSingleEmpty } from '@ui/media/icons/DotEmptySingle.tsx';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';

interface SubscriptionServiceModalProps extends PropsWithChildren {
  isOpen?: boolean;
  contractId: string;
  onClose: () => void;
  showNextInvoice?: boolean;
  billingEnabled?: boolean | null;
  allowOnlinePayment?: boolean | null;
  onChangeMode: (mode: EditModalMode) => void;
  availableBankAccount?: Store<BankAccount> | null | undefined;
}

export const ModalWithInvoicePreview = observer(
  ({
    children,
    onClose,
    isOpen,
    onChangeMode,
    showNextInvoice,
    availableBankAccount,
    billingEnabled,
    contractId,
  }: SubscriptionServiceModalProps) => {
    const store = useStore();
    const contractStore = store.contracts.value.get(contractId);

    console.log(toJS(contractStore?.value));

    const tenantBillingProfiles =
      store.settings.tenantBillingProfiles.toArray();
    const firstProfile = tenantBillingProfiles[0].value;

    const nextInvoice: TInvoice | undefined =
      contractStore?.value?.upcomingInvoices?.find((invoice: TInvoice) =>
        isAfter(new Date(invoice.issued), new Date()),
      );

    const isSimulationEnabled = useFeatureIsOn('invoice-simulation');

    useEffect(() => {
      if (isSimulationEnabled) {
        // run simulation when the edit is opened and there is no next invoice
        if (!nextInvoice && isOpen) {
          // SIMULATE INVOICE --- IS IT NEEDED NOW???
        }
        if (!isOpen) {
          // reset invoices probably not needed
        }
      }
    }, [nextInvoice, isOpen, isSimulationEnabled]);

    const billedTo = {
      addressLine1: firstProfile?.addressLine1 ?? '',
      addressLine2: firstProfile?.addressLine2 ?? '',
      locality: firstProfile.locality ?? '',
      zip: firstProfile?.zip ?? '',
      country: firstProfile?.country ?? '',
      email: firstProfile?.email ?? '',
      name: firstProfile?.legalName ?? '',
      region: firstProfile?.region ?? '',
    };
    //
    // const invoicePreviewStaticData = useMemo(
    //   () => ({
    //     status: null,
    //     invoiceNumber: 'INV-003',
    //     isBilledToFocused: false,
    //     note: '',
    //     currency: currency,
    //
    //     lines: [],
    //     tax: 0,
    //     total: 0,
    //     dueDate: new Date().toISOString(),
    //     subtotal: 0,
    //     issueDate: new Date().toISOString(),
    //
    //     from: tenantBillingProfile?.tenantBillingProfiles?.[0]
    //       ? {
    //           addressLine1:
    //             tenantBillingProfile?.tenantBillingProfiles?.[0]
    //               ?.addressLine1 ?? '',
    //           addressLine2:
    //             tenantBillingProfile?.tenantBillingProfiles?.[0].addressLine2,
    //           locality:
    //             tenantBillingProfile?.tenantBillingProfiles?.[0]?.locality ??
    //             '',
    //           zip: tenantBillingProfile?.tenantBillingProfiles?.[0]?.zip ?? '',
    //           country: tenantBillingProfile?.tenantBillingProfiles?.[0].country
    //             ? countryOptions.find(
    //                 (country) =>
    //                   country.value ===
    //                   tenantBillingProfile?.tenantBillingProfiles?.[0]?.country,
    //               )?.label
    //             : '',
    //           email:
    //             tenantBillingProfile?.tenantBillingProfiles?.[0]
    //               ?.sendInvoicesFrom,
    //           name: tenantBillingProfile?.tenantBillingProfiles?.[0]?.legalName,
    //           region:
    //             tenantBillingProfile?.tenantBillingProfiles?.[0]?.region ?? '',
    //         }
    //       : {
    //           addressLine1: '29 Maple Lane',
    //           addressLine2: 'Springfield, Haven County',
    //           locality: 'San Francisco',
    //           zip: '89302',
    //           country: 'United States of America',
    //           email: 'invoices@acme.com',
    //           name: 'Acme Corp.',
    //           region: 'California',
    //         },
    //   }),
    //   [tenantBillingProfile?.tenantBillingProfiles?.[0], currency],
    // );

    return (
      <Modal open={isOpen} onOpenChange={onClose}>
        <ModalPortal>
          <ModalOverlay className='z-50' />
          <ModalContent
            placement='center'
            className='border-r-2 flex bg-transparent shadow-none border-none z-[999] w-full gap-4'
            style={{
              minWidth: billingEnabled ? '1048px' : 'auto',
              minHeight: '80vh',
              boxShadow: 'none',
            }}
          >
            {children}

            {billingEnabled && (
              <div className='bg-white rounded-lg w-full'>
                <Invoice
                  onOpenAddressDetailsModal={() =>
                    onChangeMode(EditModalMode.BillingDetails)
                  }
                  isBilledToFocused={false}
                  shouldBlurDummy={false}
                  note={nextInvoice?.note}
                  invoiceNumber={nextInvoice?.invoiceNumber ?? ''}
                  currency={nextInvoice?.currency}
                  billedTo={billedTo}
                  from={{
                    addressLine1: nextInvoice?.provider.addressLine1 ?? '',
                    addressLine2: nextInvoice?.provider.addressLine2 ?? '',
                    locality: nextInvoice?.provider.addressLocality ?? '',
                    zip: nextInvoice?.provider.addressZip ?? '',
                    country: nextInvoice?.provider.addressCountry ?? '',
                    email: '',
                    name: nextInvoice?.provider.name ?? '',
                    region: nextInvoice?.provider.addressRegion ?? '',
                  }}
                  invoicePeriodStart={nextInvoice?.invoicePeriodStart}
                  invoicePeriodEnd={nextInvoice?.invoicePeriodEnd}
                  tax={nextInvoice?.taxDue ?? 0}
                  lines={nextInvoice?.invoiceLineItems ?? []}
                  subtotal={nextInvoice?.subtotal ?? 10}
                  issueDate={nextInvoice?.issued ?? new Date()}
                  dueDate={nextInvoice?.due ?? new Date()}
                  total={nextInvoice?.amountDue ?? 10}
                  canPayWithBankTransfer={true}
                  check={true}
                  availableBankAccount={availableBankAccount?.value}
                />
              </div>
            )}

            {/*{billingEnabled && (*/}
            {/*  <>*/}
            {/*    {showNextInvoice === undefined && (*/}
            {/*      <div*/}
            {/*        style={{ minWidth: '600px' }}*/}
            {/*        className={cn('bg-white rounded  h-full absolute z-10')}*/}
            {/*      >*/}
            {/*        <InvoiceSkeleton />*/}
            {/*      </div>*/}
            {/*    )}*/}
            {/*    {showNextInvoice && (*/}
            {/*      <div className='h-auto w-full flex relative min-w-[600px]'>*/}
            {/*        {(invoicePreviewList.simulatedInvoices.length === 0 ||*/}
            {/*          !isSimulationEnabled) && (*/}
            {/*          <div*/}
            {/*            style={{ minWidth: '600px' }}*/}
            {/*            className={cn('bg-white rounded  h-full absolute z-1 ')}*/}
            {/*          >*/}
            {/*            <div className='w-full h-full'>*/}
            {/*              <Invoice*/}
            {/*                {...invoicePreviewStaticData}*/}
            {/*                shouldBlurDummy={true}*/}
            {/*                onOpenAddressDetailsModal={() =>*/}
            {/*                  onChangeModalMode(EditModalMode.BillingDetails)*/}
            {/*                }*/}
            {/*                billedTo={billedTo}*/}
            {/*                currency={currency}*/}
            {/*                canPayWithBankTransfer={allowBankTransfer}*/}
            {/*                check={allowCheck}*/}
            {/*                availableBankAccount={availableBankAccount}*/}
            {/*              />*/}
            {/*            </div>*/}
            {/*          </div>*/}
            {/*        )}*/}

            {/*        {invoicePreviewList.simulatedInvoices.length > 1 &&*/}
            {/*          isSimulationEnabled && (*/}
            {/*            <div className='absolute top-[-30px] right-0 text-white text-base '>*/}
            {/*              <IconButton*/}
            {/*                variant='ghost'*/}
            {/*                className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'*/}
            {/*                aria-label='prev'*/}
            {/*                icon={<ChevronLeft className='text-inherit' />}*/}
            {/*                onClick={() => {*/}
            {/*                  if (*/}
            {/*                    invoicePreviewList.previewedInvoiceIndex <= 0*/}
            {/*                  ) {*/}
            {/*                    return;*/}
            {/*                  }*/}
            {/*                  invoicePreviewList.setPreviewedInvoice(*/}
            {/*                    invoicePreviewList.previewedInvoiceIndex - 1,*/}
            {/*                  );*/}
            {/*                }}*/}
            {/*              />*/}
            {/*              {invoicePreviewList.simulatedInvoices.map((e, i) => (*/}
            {/*                <IconButton*/}
            {/*                  variant='ghost'*/}
            {/*                  className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'*/}
            {/*                  key={e.invoice?.invoiceNumber}*/}
            {/*                  aria-label='prev'*/}
            {/*                  onClick={() =>*/}
            {/*                    invoicePreviewList.setPreviewedInvoice(i)*/}
            {/*                  }*/}
            {/*                  icon={*/}
            {/*                    i ===*/}
            {/*                    invoicePreviewList.previewedInvoiceIndex ? (*/}
            {/*                      <DotSingle className='text-inherit' />*/}
            {/*                    ) : (*/}
            {/*                      <DotSingleEmpty className='text-inherit' />*/}
            {/*                    )*/}
            {/*                  }*/}
            {/*                />*/}
            {/*              ))}*/}

            {/*              <IconButton*/}
            {/*                variant='ghost'*/}
            {/*                className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'*/}
            {/*                aria-label='prev'*/}
            {/*                icon={<ChevronRight className='text-inherit' />}*/}
            {/*                onClick={() => {*/}
            {/*                  if (*/}
            {/*                    invoicePreviewList.previewedInvoiceIndex >=*/}
            {/*                    invoicePreviewList.simulatedInvoices.length - 1*/}
            {/*                  ) {*/}
            {/*                    return;*/}
            {/*                  }*/}
            {/*                  invoicePreviewList.setPreviewedInvoice(*/}
            {/*                    invoicePreviewList.previewedInvoiceIndex + 1,*/}
            {/*                  );*/}
            {/*                }}*/}
            {/*              />*/}
            {/*            </div>*/}
            {/*          )}*/}

            {/*        {serviceFormStore.isSimulationRunning && (*/}
            {/*          <div*/}
            {/*            style={{ minWidth: '600px' }}*/}
            {/*            className={cn('bg-white rounded  h-full absolute z-10')}*/}
            {/*          >*/}
            {/*            <InvoiceSkeleton />*/}
            {/*          </div>*/}
            {/*        )}*/}

            {/*        {!serviceFormStore.isSimulationRunning &&*/}
            {/*          invoicePreviewList.simulatedInvoices.map(*/}
            {/*            ({ invoice }, index) => {*/}
            {/*              if (!invoice) return null;*/}

            {/*              return (*/}
            {/*                <div*/}
            {/*                  style={{ minWidth: '600px' }}*/}
            {/*                  className={cn(*/}
            {/*                    'bg-white rounded  h-full absolute z-1 ',*/}
            {/*                    {*/}
            {/*                      '-rotate-[1.15deg] shadow-lg  ':*/}
            {/*                        index !==*/}
            {/*                        invoicePreviewList.previewedInvoiceIndex,*/}
            {/*                      'shadow-md z-10 animate-zIndex':*/}
            {/*                        index ===*/}
            {/*                        invoicePreviewList.previewedInvoiceIndex,*/}
            {/*                    },*/}
            {/*                  )}*/}
            {/*                  key={`${invoice?.invoiceNumber}-${invoice.total}-invoice-preview`}*/}
            {/*                >*/}
            {/*                  <div className='w-full h-full'>*/}
            {/*                    <Invoice*/}
            {/*                      onOpenAddressDetailsModal={() =>*/}
            {/*                        onChangeModalMode(*/}
            {/*                          EditModalMode.BillingDetails,*/}
            {/*                        )*/}
            {/*                      }*/}
            {/*                      isBilledToFocused={false}*/}
            {/*                      shouldBlurDummy={false}*/}
            {/*                      note={invoice?.note}*/}
            {/*                      invoiceNumber={invoice?.invoiceNumber ?? ''}*/}
            {/*                      currency={invoice?.currency ?? currency}*/}
            {/*                      billedTo={billedTo}*/}
            {/*                      from={*/}
            {/*                        invoice*/}
            {/*                          ? {*/}
            {/*                              addressLine1:*/}
            {/*                                invoice.provider.addressLine1 ?? '',*/}
            {/*                              addressLine2:*/}
            {/*                                invoice.provider.addressLine2 ?? '',*/}
            {/*                              locality:*/}
            {/*                                invoice.provider.addressLocality ??*/}
            {/*                                '',*/}
            {/*                              zip:*/}
            {/*                                invoice.provider.addressZip ?? '',*/}
            {/*                              country:*/}
            {/*                                invoice.provider.addressCountry ??*/}
            {/*                                '',*/}
            {/*                              email: '',*/}
            {/*                              name: invoice.provider.name ?? '',*/}
            {/*                              region:*/}
            {/*                                invoice.provider.addressRegion ??*/}
            {/*                                '',*/}
            {/*                            }*/}
            {/*                          : invoicePreviewStaticData.from*/}
            {/*                      }*/}
            {/*                      invoicePeriodStart={*/}
            {/*                        invoice?.invoicePeriodStart*/}
            {/*                      }*/}
            {/*                      invoicePeriodEnd={invoice?.invoicePeriodEnd}*/}
            {/*                      tax={invoice.taxDue}*/}
            {/*                      lines={invoice?.invoiceLineItems ?? []}*/}
            {/*                      subtotal={invoice?.subtotal ?? 10}*/}
            {/*                      issueDate={invoice?.issued ?? new Date()}*/}
            {/*                      dueDate={invoice?.due ?? new Date()}*/}
            {/*                      total={invoice?.total ?? 10}*/}
            {/*                      canPayWithBankTransfer={true}*/}
            {/*                      check={allowCheck}*/}
            {/*                      availableBankAccount={availableBankAccount}*/}
            {/*                    />*/}
            {/*                  </div>*/}
            {/*                </div>*/}
            {/*              );*/}
            {/*            },*/}
            {/*          )}*/}
            {/*      </div>*/}
            {/*    )}*/}
            {/*  </>*/}
            {/*)}*/}
          </ModalContent>
        </ModalPortal>
      </Modal>
    );
  },
);
