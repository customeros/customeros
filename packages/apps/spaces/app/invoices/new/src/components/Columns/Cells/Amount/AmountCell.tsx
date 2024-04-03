import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

export const AmountCell = ({
  value,
  currency,
}: {
  value: number;
  currency: string;
}) => {
  return <span>{formatCurrency(value, 2, currency)}</span>;
};
