import { InvoiceStatus } from '@graphql/types';
import { renderStatusNode } from '@shared/components/Invoice/Cells';

export const PaymentStatusCell = ({
  value,
}: {
  value: InvoiceStatus | null;
}) => {
  return <div>{renderStatusNode(value)}</div>;
};
