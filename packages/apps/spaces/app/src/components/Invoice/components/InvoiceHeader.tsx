'use client';

import Image from 'next/image';
import React, { FC } from 'react';

import { Tag } from '@ui/presentation/Tag';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

type InvoiceHeaderProps = {
  status?: string;
  invoiceNumber: string;
};

export const InvoiceHeader: FC<InvoiceHeaderProps> = ({
  invoiceNumber,
  status,
}) => {
  const client = getGraphQLClient();

  const { data: globalCacheData } = useGlobalCacheQuery(client);

  return (
    <div>
      <div className='flex flex-1 justify-between items-center'>
        <div className='flex items-center'>
          <h1 className='text-3xl font-bold'>Invoice</h1>
          {status && (
            <div className='ml-4 mt-1'>
              <Tag variant='outline' colorScheme='gray'>
                {status}
              </Tag>
            </div>
          )}
        </div>

        {globalCacheData?.global_Cache?.cdnLogoUrl && (
          <div className='flex relative max-h-[120px] w-full justify-end'>
            <Image
              src={`${globalCacheData?.global_Cache?.cdnLogoUrl}`}
              alt='CustomerOS'
              width={136}
              height={40}
              style={{
                objectFit: 'contain',
                maxHeight: '40px',
                maxWidth: 'fit-content',
              }}
            />
          </div>
        )}
      </div>

      <h2 className='text-sm text-gray-500 '>N° {invoiceNumber}</h2>
    </div>
  );
};
