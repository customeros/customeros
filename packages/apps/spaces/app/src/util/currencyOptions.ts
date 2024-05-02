import { Currency } from '@graphql/types';

export const currencyOptions = [
  {
    label: 'EUR',
    value: Currency.Eur,
  },
  {
    label: 'GBP',
    value: Currency.Gbp,
  },
  {
    label: 'USD',
    value: Currency.Usd,
  },
];
export const currencySymbol: Record<string, string> = {
  EUR: '€',
  GBP: '£',
  USD: '$',
};
