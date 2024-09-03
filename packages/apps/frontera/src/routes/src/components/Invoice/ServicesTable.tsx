import { DateTimeUtils } from '@utils/date';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { BilledType, InvoiceLine, InvoiceLineSimulate } from '@graphql/types';

type ServicesTableProps = {
  currency: string;
  invoicePeriodEnd?: string;
  invoicePeriodStart?: string;
  billingPeriodsInMonths?: number | null;
  services: InvoiceLine[] | InvoiceLineSimulate[];
};

function isInvoiceLineSimulate(
  service: InvoiceLine | InvoiceLineSimulate,
): service is InvoiceLineSimulate {
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
  billingPeriodsInMonths,
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
              key={`invoice-line-item-${price}-${vat}-${index}-${service.description}`}
              className='flex flex-row w-full justify-between border-b border-gray-300 py-2 '
            >
              <div className={'flex w-full'}>
                <div className='w-1/2 '>
                  <div className='text-left text-sm capitalize font-medium leading-5'>
                    {service?.description ?? 'Unnamed'}
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
                                DateTimeUtils.dateWithAbreviatedMonth,
                              )}{' '}
                            {invoicePeriodEnd && invoicePeriodStart && '-'}
                            {''}
                            {invoicePeriodEnd &&
                              DateTimeUtils.format(
                                invoicePeriodEnd,
                                DateTimeUtils.dateWithAbreviatedMonth,
                              )}
                          </>
                        )}

                        {isGenerated &&
                          getBilledTypeMonths(
                            isGenerated?.contractLineItem?.billingCycle,
                          ) !== billingPeriodsInMonths && (
                            <span className='ml-2'>
                              {formatCurrency(
                                isGenerated.contractLineItem?.price,
                                2,
                                currency,
                              )}
                              {getBilledTypeLabel(
                                isGenerated.contractLineItem
                                  ?.billingCycle as BilledType,
                              )}
                            </span>
                          )}
                      </div>
                    )}
                  </div>
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {service.quantity}
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {price}
                </div>
                <div className='w-1/6 flex justify-end text-sm text-gray-500 leading-5'>
                  {vat}
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

function getBilledTypeLabel(billedType: BilledType): string {
  switch (billedType) {
    case BilledType.Annually:
      return '/year';
    case BilledType.Monthly:
      return '/month';

    case BilledType.Quarterly:
      return '/quarter';
    default:
      return '';
  }
}

function getBilledTypeMonths(billedType: BilledType): number {
  switch (billedType) {
    case BilledType.Annually:
      return 12;
    case BilledType.Monthly:
      return 1;

    case BilledType.Quarterly:
      return 3;
    default:
      return 1;
  }
}
