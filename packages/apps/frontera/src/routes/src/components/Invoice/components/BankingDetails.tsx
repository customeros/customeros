import { FC } from 'react';

import { Currency, BankAccount } from '@graphql/types';

type InvoiceHeaderProps = {
  currency?: string;
  invoiceNumber?: string | null;
  availableBankAccount?: Partial<BankAccount> | null;
};

const getBankDetails = (
  currency?: string,
  availableBankAccount?: Partial<BankAccount> | null,
): { label: string; value: string } => {
  const details = {
    label: 'BIC/Swift',
    value: availableBankAccount?.bic || '-',
  };

  switch (currency) {
    case Currency.Gbp:
      details.label = 'Sort code';
      details.value = availableBankAccount?.sortCode || '-';
      break;
    case Currency.Usd:
      details.label = 'Routing Number';
      details.value = availableBankAccount?.routingNumber || '-';
      break;
    case Currency.Eur:
      details.label = 'BIC/Swift';
      details.value = availableBankAccount?.bic || '-';
      break;
    default:
      break;
  }

  return details;
};

export const BankingDetails: FC<InvoiceHeaderProps> = ({
  availableBankAccount,
  currency,
  invoiceNumber,
}) => {
  const bankDetails: { label: string; value: string } = getBankDetails(
    currency,
    availableBankAccount,
  );

  const accountNumberLabel =
    currency === Currency.Eur ? 'IBAN' : 'Account number';
  const accountNumberValue =
    currency === Currency.Eur
      ? availableBankAccount?.iban
      : availableBankAccount?.accountNumber;

  return (
    <div className='flex flex-col border-t border-gray-300 py-2'>
      <span className='text-xs font-semibold'>Bank transfer</span>
      <div className='flex justify-between'>
        <div className='flex flex-col'>
          <span className='text-xs font-medium'>Bank name</span>
          <span className='text-xs text-gray-500'>
            {availableBankAccount?.bankName || '-'}
          </span>
        </div>
        <div className='flex flex-col'>
          <span className='text-xs font-medium'>{bankDetails.label}</span>
          <span className='text-xs text-gray-500'>{bankDetails.value}</span>
        </div>
        <div className='flex flex-col'>
          <span className='text-xs font-medium'>{accountNumberLabel}</span>
          <span className='text-xs text-gray-500'>
            {accountNumberValue || '-'}
          </span>
        </div>
        <div className='flex flex-col'>
          <span className='text-xs font-medium'>Reference</span>
          <span className='text-xs text-gray-500'>{invoiceNumber || '-'}</span>
        </div>
      </div>
    </div>
  );
};
