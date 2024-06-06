import { DateTimeUtils } from '@utils/date';
import { useStore } from '@shared/hooks/useStore';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface UpcomingInvoiceProps {
  id: string;
}

export const UpcomingInvoice = ({ id }: UpcomingInvoiceProps) => {
  const store = useStore();
  const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();
  const invoice = store.invoices.value.get(id)?.value;

  return (
    <div
      key={invoice?.metadata?.id}
      className='flex  text-sm'
      role='button'
      tabIndex={0}
      onClick={() => handleOpenInvoice(invoice?.metadata?.id)}
    >
      <div className='whitespace-nowrap mr-1'>Monthly recurring:</div>
      <div className='whitespace-nowrap text-gray-500 underline'>
        {formatCurrency(invoice.amountDue, 2, invoice?.currency)} on{' '}
        {DateTimeUtils.format(
          invoice?.due,
          DateTimeUtils.defaultFormatShortString,
        )}{' '}
        (
        {DateTimeUtils.format(
          invoice?.invoicePeriodStart,
          DateTimeUtils.dateDayAndMonth,
        )}{' '}
        -{' '}
        {DateTimeUtils.format(
          invoice?.invoicePeriodEnd,
          DateTimeUtils.dateDayAndMonth,
        )}
        )
      </div>
    </div>
  );
};
