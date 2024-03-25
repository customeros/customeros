import React, { ReactNode, ReactElement } from 'react';

import { Currency } from '@graphql/types';
import { Text } from '@ui/typography/Text';

const CurrencySymbol = ({ children }: { children: ReactNode }) => {
  return (
    <Text
      textAlign='center'
      color='gray.500'
      whiteSpace='nowrap'
      fontWeight='semibold'
      fontSize='sm'
    >
      {children}
    </Text>
  );
};
export const currencyIcon: Record<string, ReactElement> = {
  [Currency.Usd]: <CurrencySymbol>$</CurrencySymbol>,
  [Currency.Gbp]: <CurrencySymbol>€</CurrencySymbol>,
  [Currency.Eur]: <CurrencySymbol>£</CurrencySymbol>,
  [Currency.Aud]: <CurrencySymbol>AU$</CurrencySymbol>,
  [Currency.Brl]: <CurrencySymbol>R$</CurrencySymbol>,
  [Currency.Cad]: <CurrencySymbol>CA$</CurrencySymbol>,
  [Currency.Chf]: <CurrencySymbol>Fr.</CurrencySymbol>,
  [Currency.Cny]: <CurrencySymbol>CN¥</CurrencySymbol>,
  [Currency.Hkd]: <CurrencySymbol>HK$</CurrencySymbol>,
  [Currency.Inr]: <CurrencySymbol>Rs.</CurrencySymbol>,
  [Currency.Jpy]: <CurrencySymbol>¥</CurrencySymbol>,
  [Currency.Krw]: <CurrencySymbol>₩</CurrencySymbol>,
  [Currency.Mxn]: <CurrencySymbol>MX$</CurrencySymbol>,
  [Currency.Nok]: <CurrencySymbol>kr</CurrencySymbol>,
  [Currency.Nzd]: <CurrencySymbol>NZ$</CurrencySymbol>,
  [Currency.Ron]: <CurrencySymbol>L</CurrencySymbol>,
  [Currency.Sek]: <CurrencySymbol>kr</CurrencySymbol>,
  [Currency.Sgd]: <CurrencySymbol>S$</CurrencySymbol>,
  [Currency.Try]: <CurrencySymbol>TL</CurrencySymbol>,
  [Currency.Zar]: <CurrencySymbol>R</CurrencySymbol>,
};
export function mapCurrencyToOptions() {
  return Object.values(Currency)
    .map((key) => ({
      label: key,
      value: key,
    }))
    .filter(
      (e) => ![Currency.Eur, Currency.Usd, Currency.Gbp].includes(e.value),
    );
}

export const paymentMethods: Record<string, string> = {
  card: 'Debit or Card',
  ach_debit: 'ACH',
  sepa: 'SEPA',
  bacs_debit: 'Bacs',
};
