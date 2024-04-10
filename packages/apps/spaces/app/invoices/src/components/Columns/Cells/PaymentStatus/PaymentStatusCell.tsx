import { PaymentStatusSelect } from '@invoices/components/shared';

import { InvoiceStatus } from '@graphql/types';

export const PaymentStatusCell = ({
  value,
  invoiceId,
}: {
  invoiceId: string;
  value: InvoiceStatus | null;
}) => {
  return <PaymentStatusSelect value={value} invoiceId={invoiceId} />;
};
