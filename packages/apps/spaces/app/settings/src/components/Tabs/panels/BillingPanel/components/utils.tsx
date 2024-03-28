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
  [Currency.Gbp]: <CurrencySymbol>£</CurrencySymbol>,
  [Currency.Eur]: <CurrencySymbol>€</CurrencySymbol>,
};

export const paymentMethods: Record<string, string> = {
  card: 'Credit or Debit card',
  ach_debit: 'ACH',
  sepa: 'SEPA',
  bacs_debit: 'Bacs',
};
