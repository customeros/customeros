import { useRef, useState, useEffect, PropsWithChildren } from 'react';

import { autorun } from 'mobx';
import { Store } from '@store/store.ts';
import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { Invoice as TInvoice } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { countryOptions } from '@shared/util/countryOptions';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { BankAccount, InvoiceSimulate } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient.ts';
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
  contractId: string;
  showNextInvoice?: boolean;
  billingEnabled?: boolean | null;
  allowOnlinePayment?: boolean | null;
  availableBankAccount?: Store<BankAccount> | null | undefined;
}
type TSimulateQueryResult = { invoice_Simulate: Array<InvoiceSimulate> };

export const ModalWithInvoicePreview = observer(
  ({
    children,
    showNextInvoice,
    availableBankAccount,
    billingEnabled,
    contractId,
  }: SubscriptionServiceModalProps) => {
    const abortControllerRef = useRef<AbortController | null>(null);
    const client = getGraphQLClient({
      signal: abortControllerRef?.current?.signal,
    });

    const { isEditModalOpen, onChangeModalMode, onEditModalClose } =
      useContractModalStateContext();
    const store = useStore();
    const [isSimulatingInvoices, setIsSimulatingInvoices] = useState(false);
    const [previewedInvoiceIndex, setPreviewedInvoiceIndex] = useState(0);
    const [simulatedInvoices, setSimulatedInvoices] = useState<
      InvoiceSimulate[]
    >([]);

    const contractStore = store.contracts.value.get(contractId);
    const tenantBillingProfiles =
      store.settings.tenantBillingProfiles.toArray();

    const nextInvoice: TInvoice | undefined = store.invoices
      .toArray()
      ?.find(
        (e) =>
          contractStore?.value?.upcomingInvoices?.[0]?.metadata?.id ===
          e.value?.metadata?.id,
      )?.value;

    const isSimulationEnabled = useFeatureIsOn('invoice-simulation');

    const runSimulation = async (): Promise<void> => {
      const payload = contractStore?.value.contractLineItems
        ?.map((e) =>
          store.contractLineItems
            ?.toArray()
            .filter((d) => d.value.metadata.id === e.metadata.id),
        )
        .flat()
        .map((line) => {
          return {
            key: line?.value?.metadata?.id,
            serviceLineItemId: line?.value?.metadata?.id,
            parentId: line?.value?.parentId,
            description: line?.value?.description,
            billingCycle: line?.value?.billingCycle,
            price: line?.value?.price,
            quantity: line?.value?.quantity,
            serviceStarted: line?.value?.serviceStarted,
            taxRate: line?.value?.tax?.taxRate,
            closeVersion: line?.value?.closed,
          };
        });

      setIsSimulatingInvoices(true);

      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
      const abortController = new AbortController();

      abortControllerRef.current = abortController;

      try {
        const { invoice_Simulate } = await client.request<TSimulateQueryResult>(
          SimulateInvoiceDocument,
          {
            input: {
              contractId,
              serviceLines: payload,
            },
          },
        );
        setSimulatedInvoices(invoice_Simulate);
      } catch (error) {
        console.error(`Simulation failed: ${error}`);
      } finally {
        setIsSimulatingInvoices(false);
      }
    };

    useEffect(() => {
      if (isSimulationEnabled && !nextInvoice && isEditModalOpen) {
        const disposer = autorun(
          () => {
            if (store.contractLineItems) {
              runSimulation();
            }
          },
          {
            delay: 1000,
            requiresObservable: true,
          },
        );

        return () => disposer();
      }
    }, [nextInvoice, isEditModalOpen, isSimulationEnabled]);

    useEffect(() => {
      return () => {
        if (abortControllerRef.current) {
          abortControllerRef.current.abort();
        }
      };
    }, []);

    const billedTo = {
      addressLine1: contractStore?.value?.billingDetails?.addressLine1 ?? '',
      addressLine2: contractStore?.value?.billingDetails?.addressLine2 ?? '',
      locality: contractStore?.value?.billingDetails?.locality ?? '',
      zip: contractStore?.value?.billingDetails?.postalCode ?? '',
      country: contractStore?.value?.billingDetails?.country ?? '',
      email: contractStore?.value?.billingDetails?.billingEmail ?? '',
      name: contractStore?.value?.billingDetails?.organizationLegalName ?? '',
      region: contractStore?.value?.billingDetails?.region ?? '',
    };

    const invoicePreviewStaticData = {
      status: null,
      invoiceNumber: 'INV-003',
      isBilledToFocused: false,
      note: '',
      currency: contractStore?.value?.currency ?? 'USD',

      lines: [],
      tax: 0,
      total: 0,
      dueDate: new Date().toISOString(),
      subtotal: 0,
      issueDate: new Date().toISOString(),

      from: tenantBillingProfiles?.[0]
        ? {
            addressLine1: tenantBillingProfiles?.[0]?.value?.addressLine1 ?? '',
            addressLine2: tenantBillingProfiles?.[0]?.value?.addressLine2,
            locality: tenantBillingProfiles?.[0]?.value?.locality ?? '',
            zip: tenantBillingProfiles?.[0]?.value?.zip ?? '',
            country: tenantBillingProfiles?.[0]?.value?.country
              ? countryOptions.find(
                  (country) =>
                    country.value ===
                    tenantBillingProfiles?.[0]?.value?.country,
                )?.label
              : '',
            email: tenantBillingProfiles?.[0]?.value?.sendInvoicesFrom,
            name: tenantBillingProfiles?.[0]?.value?.legalName,
            region: tenantBillingProfiles?.[0]?.value?.region ?? '',
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
    };

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
                    {(simulatedInvoices.length === 0 ||
                      !isSimulationEnabled) && (
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
                            currency={contractStore?.value.currency ?? 'USD'}
                            canPayWithBankTransfer={
                              contractStore?.value?.billingDetails
                                ?.canPayWithBankTransfer
                            }
                            check={contractStore?.value?.billingDetails?.check}
                            availableBankAccount={availableBankAccount?.value}
                          />
                        </div>
                      </div>
                    )}

                    {simulatedInvoices.length > 1 && isSimulationEnabled && (
                      <div className='absolute top-[-30px] right-0 text-white text-base '>
                        <IconButton
                          variant='ghost'
                          className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'
                          aria-label='prev'
                          icon={<ChevronLeft className='text-inherit' />}
                          onClick={() => {
                            if (previewedInvoiceIndex <= 0) {
                              return;
                            }
                            setPreviewedInvoiceIndex(previewedInvoiceIndex - 1);
                          }}
                        />
                        {simulatedInvoices.map((e, i) => (
                          <IconButton
                            variant='ghost'
                            className='bg-transparent text-white p-0 hover:text-white hover:bg-transparent focus:text-white focus:bg-transparent'
                            key={e?.invoiceNumber}
                            aria-label='prev'
                            onClick={() => setPreviewedInvoiceIndex(i)}
                            icon={
                              i === previewedInvoiceIndex ? (
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
                              previewedInvoiceIndex >=
                              simulatedInvoices.length - 1
                            ) {
                              return;
                            }
                            setPreviewedInvoiceIndex(previewedInvoiceIndex + 1);
                          }}
                        />
                      </div>
                    )}

                    {isSimulatingInvoices && (
                      <div
                        style={{ minWidth: '600px' }}
                        className={cn('bg-white rounded  h-full absolute z-10')}
                      >
                        <InvoiceSkeleton />
                      </div>
                    )}

                    {!isSimulatingInvoices &&
                      simulatedInvoices.length &&
                      simulatedInvoices.map((invoice, index) => {
                        if (!invoice) return null;

                        return (
                          <div
                            style={{ minWidth: '600px' }}
                            className={cn(
                              'bg-white rounded  h-full absolute z-1 ',
                              {
                                '-rotate-[1.15deg] shadow-lg  ':
                                  index !== previewedInvoiceIndex,
                                'shadow-md z-10 animate-zIndex':
                                  index === previewedInvoiceIndex,
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
                                currency={
                                  invoice?.currency ??
                                  contractStore?.value?.currency
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
                                          invoice.provider.addressLocality ??
                                          '',
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
                                check={
                                  contractStore?.value?.billingDetails?.check ??
                                  false
                                }
                                availableBankAccount={
                                  availableBankAccount?.value
                                }
                              />
                            </div>
                          </div>
                        );
                      })}
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

const SimulateInvoiceDocument = `
mutation simulateInvoice($input: InvoiceSimulateInput!) {
  invoice_Simulate(input: $input) {
    amount
    currency
    due
    invoiceNumber
    invoicePeriodEnd
    invoicePeriodStart
    issued
    note
    offCycle
    postpaid
    subtotal
    taxDue
    total
    invoiceLineItems {
      key
      description
      price
      quantity
      subtotal
      taxDue
      total
    }
    customer {
      name
      email
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
      addressRegion
    }
    provider {
      name
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
      addressRegion
    }
  }
}
    `;
