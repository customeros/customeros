import React, { ReactNode, ReactElement } from 'react';

import { Currency } from '@graphql/types';

const CurrencySymbol = ({ children }: { children: ReactNode }) => {
  return (
    <div className='align-middle text-gray-500 whitespace-nowrap font-semibold text-sm'>
      {children}
    </div>
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
