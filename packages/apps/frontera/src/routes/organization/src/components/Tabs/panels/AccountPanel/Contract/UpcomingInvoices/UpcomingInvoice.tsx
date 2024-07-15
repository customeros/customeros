import { observer } from 'mobx-react-lite';

import { Contract } from '@graphql/types';
import { DateTimeUtils } from '@utils/date';
import { useStore } from '@shared/hooks/useStore';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

function getCommittedPeriodLabel(months: string | number) {
  if (`${months}` === '1') {
    return 'Monthly';
  }
  if (`${months}` === '3') {
    return 'Quarterly';
  }

  if (`${months}` === '12') {
    return 'Annual';
  }

  return `${months}-month`;
}

interface UpcomingInvoiceProps {
  id: string;
  contractId: string;
}

export const UpcomingInvoice = observer(
  ({ id, contractId }: UpcomingInvoiceProps) => {
    const store = useStore();
    const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();
    const invoice = store.invoices.value.get(id)?.value;

    if (!invoice?.metadata.id) return null;
    const contract = store.contracts.value.get(contractId)?.value as Contract;
    const renewalPeriod = getCommittedPeriodLabel(
      contract?.committedPeriodInMonths,
    );
    const autoRenewal = contract?.autoRenew;

    return (
      <div
        key={invoice.metadata.id}
        className='flex  text-sm'
        role='button'
        tabIndex={0}
        onClick={() => handleOpenInvoice(invoice.metadata.id)}
      >
        <div className='whitespace-nowrap mr-1'>
          {renewalPeriod} {autoRenewal ? 'recurring' : ''}:
        </div>
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
            DateTimeUtils.dateWithShortYear,
          )}
          )
        </div>
      </div>
    );
  },
);
