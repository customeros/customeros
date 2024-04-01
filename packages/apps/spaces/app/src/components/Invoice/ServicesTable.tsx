import React from 'react';

import { Invoice } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

type ServicesTableProps = {
  currency: string;
  services: Partial<Invoice['invoiceLineItems']>;
};

export function ServicesTable({ services, currency }: ServicesTableProps) {
  return (
    <div className='w-full'>
      <div className='flex flex-col w-full'>
        <div className='flex flex-row w-full justify-between border-b border-gray-300 py-2'>
          <div className='w-1/2 text-left text-sm capitalize font-bold'>
            Service
          </div>
          <div className='w-1/6 text-center text-sm capitalize font-bold'>
            Qty
          </div>
          <div className='w-1/6 text-center text-sm capitalize font-bold'>
            Unit Price
          </div>
          <div className='w-1/6 text-right text-sm capitalize font-bold'>
            Amount
          </div>
        </div>
        {services.map((service, index) => (
          <div
            className='flex flex-row w-full justify-between border-b border-gray-300 py-4'
            key={index}
          >
            <div className='w-1/2 text-left text-sm capitalize font-medium leading-5'>
              {service?.description ?? 'Unnamed'}
            </div>
            <div className='w-1/6 text-center text-sm text-gray-500 leading-5'>
              {service?.quantity}
            </div>
            <div className='w-1/6 text-center text-sm text-gray-500 leading-5'>
              {formatCurrency(service?.price ?? 0, 2, currency)}
            </div>
            <div className='w-1/6 text-right text-sm text-gray-500 leading-5'>
              {formatCurrency(service?.total ?? 0, 2, currency)}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
