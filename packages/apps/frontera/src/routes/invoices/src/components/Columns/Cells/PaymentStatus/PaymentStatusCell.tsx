import { InvoiceStatus } from '@graphql/types';

import { PaymentStatusSelect } from '../../../shared';

export const PaymentStatusCell = ({
  value,
  invoiceId,
}: {
  invoiceId: string;
  value: InvoiceStatus | null;
}) => {
  return <PaymentStatusSelect value={value} invoiceId={invoiceId} />;
};
