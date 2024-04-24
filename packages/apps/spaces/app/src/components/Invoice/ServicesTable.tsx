import React from 'react';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@spaces/utils/date';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { BilledType, InvoiceLine, InvoiceLineSimulate } from '@graphql/types';
import { Highlighter } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/highlighters';
import { ISimulatedInvoiceLineItems } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/InvoicePreviewList.store';

type ServicesTableProps = {
  currency: string;
  invoicePeriodEnd?: string;
  shouldBlurDummy?: boolean;
  invoicePeriodStart?: string;
  services: InvoiceLine[] | ISimulatedInvoiceLineItems[];
};
function isInvoiceLineSimulate(
  service: InvoiceLine | (InvoiceLineSimulate & { color: string }),
): service is InvoiceLineSimulate & { color: string } {
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
  shouldBlurDummy,
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
            Amount
          </div>
          <div className='w-1/6 text-right text-sm capitalize font-bold'>
            VAT
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
              <div
                className={cn('flex w-full', {
                  'filter-none': !shouldBlurDummy,
                  'blur-[2px]': shouldBlurDummy,
                })}
              >
                <div className='w-1/2 '>
                  <div className='text-left text-sm capitalize font-medium leading-5'>
                    {isGenerated && (isGenerated?.description ?? 'Unnamed')}
                    {isSimulated && (
                      <div className='max-w-fit'>
                        {isSimulated.description ?? 'Unnamed'}
                      </div>
                    )}
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
                          <Highlighter
                            backgroundColor={
                              isSimulated.diff.find((e) => e === 'billingCycle')
                                ? isSimulated.color
                                : 'transparent'
                            }
                            highlightVersion={isSimulated.shapeVariant}
                          >
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
                          </Highlighter>
                        )}
                      </div>
                    )}
                  </div>
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {isSimulated ? (
                    <div className='max-w-fit'>
                      <Highlighter
                        highlightVersion={isSimulated.shapeVariant}
                        backgroundColor={
                          isSimulated.diff.find((e) => e === 'quantity')
                            ? isSimulated.color
                            : 'transparent'
                        }
                      >
                        {isSimulated.quantity ?? 'Unnamed'}
                      </Highlighter>
                    </div>
                  ) : (
                    service.quantity
                  )}
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {isSimulated ? (
                    <div className='max-w-fit'>
                      <Highlighter
                        highlightVersion={isSimulated.shapeVariant}
                        backgroundColor={
                          isSimulated.diff.find((e) => e === 'price')
                            ? isSimulated.color
                            : 'transparent'
                        }
                      >
                        {price}
                      </Highlighter>
                    </div>
                  ) : (
                    price
                  )}
                </div>
                <div className='w-1/6 text-right text-sm text-gray-500 leading-5'>
                  {formatCurrency(service?.total ?? 0, 2, currency)}
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {isSimulated ? (
                    <div className='max-w-fit'>
                      <Highlighter
                        highlightVersion={isSimulated.shapeVariant}
                        backgroundColor={
                          isSimulated.diff.find((e) => e === 'taxRate')
                            ? isSimulated.color
                            : 'transparent'
                        }
                      >
                        {vat}
                      </Highlighter>
                    </div>
                  ) : (
                    vat
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
