import React, { useState } from 'react';

import { ChevronRight } from '@ui/media/icons/ChevronRight';

import { CheckoutCard } from './components/CheckoutCard';
import { EmptyMailboxes } from './components/EmptyMailboxes';
import { AddDomainsCard } from './components/AddDomainsCard';
import { UsersCard } from './components/UsersCard/UsersCard';

export const Mailboxes = () => {
  const [isUpdated, setIsUpdated] = useState(true);

  const handleUpdate = () => {
    setIsUpdated(true);
  };

  return (
    <>
      {!isUpdated ? (
        <EmptyMailboxes onUpdate={handleUpdate} />
      ) : (
        <div className='grid grid-cols-2 gap-2 max-w-[800px] border-r-[1px] h-full'>
          <div className='py-[10px] px-6 flex flex-col border-r-[1px] h-full'>
            <div className='flex items-center justify-start gap-1 mb-4'>
              <span className='text-gray-500 font-semibold'>Mailboxes</span>
              <ChevronRight className='mt-0.5 text-gray-400 size-3' />
              <span className='font-semibold'>Add new</span>
            </div>
            <div className=' space-y-4 '>
              <AddDomainsCard />
              <UsersCard />
            </div>
          </div>
          <div className='py-[10px] px-6 flex flex-col h-full'>
            <p className='mb-4 font-semibold'>Checkout</p>
            <CheckoutCard />
          </div>
        </div>
      )}
    </>
  );
};
