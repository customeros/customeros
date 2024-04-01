'use client';

import React from 'react';

import { Skeleton } from '@ui/feedback/Skeleton';
import { Divider } from '@ui/presentation/Divider/Divider';

import { ServicesTable } from './ServicesTable';

export function InvoiceSkeleton() {
  return (
    <div className='flex flex-col px-4'>
      <div className='flex flex-col mt-2'>
        <div className='flex items-center'>
          <h1 className='text-3xl font-bold'>Invoice</h1>
        </div>

        <div className='flex items-center text-gray-500 text-sm'>
          N° <Skeleton className='w-[60px] h-3 ml-1' />
        </div>

        <div className='flex mt-1 justify-evenly gap-3 border-t border-b border-gray-300'>
          <div className='flex flex-col flex-1 min-w-[150px] py-2 border-r border-gray-300 ml-2'>
            <span className='font-semibold mb-2 text-sm'>Issued</span>
            <Skeleton className='w-[70px] h-3' />
            <span className='font-semibold mt-5 mb-2 text-sm'>Due</span>
            <Skeleton className='w-[70px] h-3' />
          </div>
          <div className='flex flex-col flex-1 min-w-[150px] py-2 border-gray-300 relative'>
            <span className='font-semibold mb-2 text-sm'>Billed to</span>
            <Skeleton className='w-[90px] h-3 mb-2' />
            <Skeleton className='w-[110px] h-3 mb-1' />
            <Skeleton className='w-[50px] h-3 mb-1' />

            <div className='flex'>
              <Skeleton className='w-15 h-3 mr-2 mb-1' />
              <Skeleton className='w-10 h-3 mb-1' />
            </div>
            <Skeleton className='w-10 h-3 mb-2' />
            <Skeleton className='w-[90px] h-3' />
          </div>
          <div className='flex flex-col flex-1 min-w-[150px] py-2'>
            <span className='font-semibold mb-2 text-sm'>From</span>
            <Skeleton className='w-25 h-3 mb-2' />
            <Skeleton className='w-30 h-3 mb-1' />
            <Skeleton className='w-[50px] h-3 mb-1' />

            <div className='flex'>
              <Skeleton className='w-15 h-3 mr-2 mb-1' />
              <Skeleton className='w-10 h-3 mb-1' />
            </div>
            <Skeleton className='w-10 h-3 mb-1' />
            <Skeleton className='w-[90px] h-3' />
          </div>
        </div>
      </div>

      <div className='flex flex-col mt-4'>
        <ServicesTable services={[]} currency='USD' />
        <div className='flex my-5 justify-between'>
          <Skeleton className='w-[55%] h-[14px] mr-2' />
          <Skeleton className='w-[10%] h-[14px] mr-2' />
          <Skeleton className='w-[20%] h-[14px] mr-2' />
          <Skeleton className='w-[15%] h-[14px] mr-2' />
        </div>

        <div className='flex flex-col self-end w-[50%] max-w-[300px] mt-4'>
          <div className='flex justify-between items-center'>
            <span className='text-sm font-medium'>Subtotal</span>
            <Skeleton className='w-5 h-3' />
          </div>
          <Divider className='my-1' />
          <div className='flex justify-between items-center'>
            <span className='text-sm'>Tax</span>
            <Skeleton className='w-5 h-3' />
          </div>
          <Divider className='my-1' />
          <div className='flex justify-between items-center'>
            <span className='text-sm font-medium'>Total</span>
            <Skeleton className='w-5 h-3' />
          </div>
          <Divider className='my-1 border-gray-500' />
          <div className='flex justify-between items-center'>
            <span className='text-sm font-semibold'>Amount due</span>
            <Skeleton className='w-5 h-3' />
          </div>
          <Divider className='my-1 border-gray-500' />
        </div>
      </div>
    </div>
  );
}
