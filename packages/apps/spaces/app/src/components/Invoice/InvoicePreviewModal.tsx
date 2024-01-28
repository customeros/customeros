import React from 'react';

import { Text } from '@ui/typography/Text';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { GetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
interface InvoicePreviewModalProps {
  id: string;
  isError: boolean;
  isFetching: boolean;
  data: GetInvoiceQuery | undefined;
}

export const InvoicePreviewModalContent: React.FC<InvoicePreviewModalProps> = ({
  id,
  isFetching,
  isError,
  data,
}) => {
  if (isFetching) {
    return <InvoiceSkeleton />;
  }

  if (!data?.invoice || isError) {
    // eslint-disable-next-line react/no-unescaped-entities
    return <Text> Sorry, we couldn't find this invoice</Text>;
  }

  return (
    <Invoice
      tax={data?.invoice?.vat}
      note={''}
      from={{
        zip: '',
        email: '',
        name: '',
        country: '',
        locality: '',
        addressLine: '',
        addressLine2: '',
      }}
      total={data?.invoice.totalAmount}
      dueDate={data?.invoice.dueDate}
      subtotal={data?.invoice.amount}
      issueDate={data?.invoice?.createdAt}
      billedTo={{
        zip: '',
        email: '',
        name: '',
        country: '',
        locality: '',
        addressLine: '',
        addressLine2: '',
      }}
      invoiceNumber={data?.invoice?.number ?? ''}
      lines={data?.invoice?.invoiceLines ?? []}
      currency={data?.invoice?.currency ?? 'USD'}
    />
  );
};
