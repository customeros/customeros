import { BilledType } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';

export function formatString(str: string, type: string, currency: string) {
  const digitCount = type === BilledType.Usage ? 4 : 2;
  const regex =
    type === BilledType.Usage ? /\b(\d+\.\d{4})\b/g : /\b(\d+\.\d{2})\b/g;

  return str.replace(regex, (_, number) => {
    return formatCurrency(Number(number), digitCount, currency);
  });
}
