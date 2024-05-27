import { DateTimeUtils } from '@spaces/utils/date';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { BilledType, InvoiceLine, InvoiceLineSimulate } from '@graphql/types';
import { Highlighter } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/Services/components/highlighters';
import { ISimulatedInvoiceLineItems } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/InvoicePreviewList.store.ts';

type ServicesTableProps = {
  currency: string;
  invoicePeriodEnd?: string;
  invoicePeriodStart?: string;
  services: InvoiceLine[] | ISimulatedInvoiceLineItems[];
};

function isInvoiceLineSimulate(
  service: InvoiceLine | ISimulatedInvoiceLineItems,
): service is ISimulatedInvoiceLineItems {
  return service && (service as InvoiceLineSimulate).key !== null;
}

function isPartialInvoiceLineItem(
  service: InvoiceLine | InvoiceLineSimulate,
): service is InvoiceLine {
  return service && (service as InvoiceLine).contractLineItem !== undefined;
}

export function ServicesTable({
  services,
  currency,
  invoicePeriodStart,
  invoicePeriodEnd,
}: ServicesTableProps) {
  return (
    <div className='w-full'>
      <div className='flex flex-col w-full'>
        <div className='flex flex-row w-full justify-between border-b border-gray-300 py-2'>
          <div className='w-1/2 text-left text-sm capitalize font-bold'>
            Service
          </div>
          <div className='w-1/6 text-right text-sm capitalize font-bold'>
            Qty
          </div>
          <div className='w-1/6 text-right text-sm capitalize font-bold'>
            Unit Price
          </div>
          <div className='w-1/6 text-right text-sm capitalize font-bold'>
            VAT
          </div>
          <div className='w-1/6 text-right text-sm capitalize font-bold'>
            Amount
          </div>
        </div>
        {services?.map((service, index) => {
          const isSimulated = isInvoiceLineSimulate(service) && service;
          const isGenerated = isPartialInvoiceLineItem(service) && service;
          const price = formatCurrency(service?.price ?? 0, 2, currency);
          const vat = formatCurrency(service?.taxDue ?? 0, 2, currency);

          return (
            <div
              className='flex flex-row w-full justify-between border-b border-gray-300 py-2 '
              key={`invoice-line-item-${price}-${vat}-${index}-${service.description}`}
            >
              <div className={'flex w-full'}>
                <div className='w-1/2 '>
                  <div className='text-left text-sm capitalize font-medium leading-5'>
                    {isGenerated && (isGenerated?.description ?? 'Unnamed')}
                  </div>
                  <div className='text-gray-500 text-sm'>
                    {isGenerated &&
                    isGenerated?.contractLineItem?.billingCycle ===
                      BilledType.Once ? (
                      <>
                        {service?.contractLineItem?.serviceStarted &&
                          DateTimeUtils.format(
                            service.contractLineItem.serviceStarted,
                            DateTimeUtils.defaultFormatShortString,
                          )}
                      </>
                    ) : (
                      <div className='max-w-fit'>
                        {isSimulated && (
                          <>
                            {invoicePeriodStart &&
                              DateTimeUtils.format(
                                invoicePeriodStart,
                                DateTimeUtils.defaultFormatShortString,
                              )}{' '}
                            {invoicePeriodEnd && invoicePeriodStart && '-'}
                            {''}
                            {invoicePeriodEnd &&
                              DateTimeUtils.format(
                                invoicePeriodEnd,
                                DateTimeUtils.defaultFormatShortString,
                              )}
                          </>
                        )}
                      </div>
                    )}
                  </div>
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {isSimulated &&
                  isSimulated.serviceLineItemStore?.isFieldRevised(
                    'quantity',
                  ) ? (
                    <div className='max-w-fit'>
                      <Highlighter
                        highlightVersion={
                          isSimulated.serviceLineItemStore.uiMetadata
                            ?.shapeVariant
                        }
                        backgroundColor={
                          isSimulated.serviceLineItemStore.uiMetadata?.color
                        }
                      >
                        {isSimulated.quantity ?? '0'}
                      </Highlighter>
                    </div>
                  ) : (
                    service.quantity
                  )}
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {isSimulated &&
                  isSimulated.serviceLineItemStore?.isFieldRevised('price') ? (
                    <div className='max-w-fit'>
                      <Highlighter
                        highlightVersion={
                          isSimulated.serviceLineItemStore.uiMetadata
                            ?.shapeVariant
                        }
                        backgroundColor={
                          isSimulated.serviceLineItemStore.uiMetadata?.color
                        }
                      >
                        {price}
                      </Highlighter>
                    </div>
                  ) : (
                    price
                  )}
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {isSimulated &&
                  isSimulated.serviceLineItemStore?.isFieldRevised(
                    'taxRate',
                  ) ? (
                    <div className='max-w-fit'>
                      <Highlighter
                        highlightVersion={
                          isSimulated.serviceLineItemStore.uiMetadata
                            ?.shapeVariant
                        }
                        backgroundColor={
                          isSimulated.serviceLineItemStore.uiMetadata?.color
                        }
                      >
                        {vat}
                      </Highlighter>
                    </div>
                  ) : (
                    vat
                  )}
                </div>
                <div className='w-1/6 text-right text-sm text-gray-500 leading-5'>
                  {formatCurrency(service?.total ?? 0, 2, currency)}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
