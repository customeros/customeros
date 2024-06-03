import { FC } from 'react';

import { cn } from '@ui/utils/cn';
import { Divider } from '@ui/presentation/Divider/Divider';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

interface InvoiceSummaryProps {
  tax: number;
  total: number;
  subtotal: number;
  currency: string;
  amountDue?: number;
  note?: string | null;
  shouldBlurDummy?: boolean;
}

export const InvoiceSummary: FC<InvoiceSummaryProps> = ({
  subtotal,
  tax,
  total,
  amountDue,
  currency,
  note,
  shouldBlurDummy,
}) => {
  return (
    <div className='flex flex-col self-end w-[50%] max-w-[300px] mt-4'>
      <div className='flex justify-between'>
        <span className='text-sm items-center font-medium'>Subtotal</span>
        <span
          className={cn('text-sm ml-2 span-gray-600', {
            'filter-none': !shouldBlurDummy,
            'blur-[2px]': shouldBlurDummy,
          })}
        >
          {formatCurrency(subtotal, 2, currency)}
        </span>
      </div>
      <Divider className='my-1 border-gray-300' />

      <div className='flex justify-between'>
        <span className='text-sm'>Tax</span>
        <span
          className={cn('text-sm ml-2 span-gray-600', {
            'filter-none': !shouldBlurDummy,
            'blur-[2px]': shouldBlurDummy,
          })}
        >
          {formatCurrency(tax, 2, currency)}
        </span>
      </div>
      <Divider className='my-1 border-gray-300' />

      <div className='flex justify-between'>
        <span className='text-sm font-medium'>Total</span>
        <span
          className={cn('text-sm ml-2 text-gray-600', {
            'filter-none': !shouldBlurDummy,
            'blur-[2px]': shouldBlurDummy,
          })}
        >
          {formatCurrency(total, 2, currency)}
        </span>
      </div>
      <Divider className='my-1 border-gray-500' />

      <div className='flex justify-between'>
        <span className='text-sm font-medium'>Amount due</span>
        <span
          className={cn('text-sm font-semibold ml-2', {
            'filter-none': !shouldBlurDummy,
            'blur-[2px]': shouldBlurDummy,
          })}
        >
          {formatCurrency(amountDue || total, 2, currency)}
        </span>
      </div>
      <Divider className='my-1 border-gray-500' />

      {note && (
        <div>
          <span className='text-sm font-medium'>Note:</span>
          <span className='text-sm ml-2 text-gray-500'>{note}</span>
        </div>
      )}
    </div>
  );
};
