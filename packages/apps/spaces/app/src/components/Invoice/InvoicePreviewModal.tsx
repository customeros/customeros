import React from 'react';

import { Text } from '@ui/typography/Text';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { InvoiceCustomer, InvoiceProvider } from '@graphql/types';
import { GetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
interface InvoicePreviewModalProps {
  id: string;
  isError: boolean;
  isFetching: boolean;
  data: GetInvoiceQuery | undefined;
}

const extractAddressData = (invoiceData: InvoiceCustomer | InvoiceProvider) => {
  return {
    zip: invoiceData?.addressZip ?? '',
    email: (invoiceData as InvoiceCustomer)?.email ?? '',
    name: invoiceData?.name ?? '',
    country: invoiceData?.addressCountry ?? '',
    locality: invoiceData?.addressLocality ?? '',
    addressLine: invoiceData?.addressLine1 ?? '',
    addressLine2: invoiceData?.addressLine2 ?? '',
  };
};

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

  const customerAddressData = extractAddressData(data?.invoice?.customer);
  const providerAddressData = extractAddressData(data?.invoice?.provider);

  return (
    <Invoice
      tax={data?.invoice?.vat}
      note={''}
      from={providerAddressData}
      total={data?.invoice.totalAmount}
      dueDate={data?.invoice.dueDate}
      subtotal={data?.invoice.amount}
      issueDate={data?.invoice?.createdAt}
      billedTo={customerAddressData}
      invoiceNumber={data?.invoice?.number ?? ''}
      lines={data?.invoice?.invoiceLines ?? []}
      currency={data?.invoice?.currency ?? 'USD'}
    />
  );
};
