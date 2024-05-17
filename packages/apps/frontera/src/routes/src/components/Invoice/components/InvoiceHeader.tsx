import { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { Tag } from '@ui/presentation/Tag';
import { Image } from '@ui/media/Image/Image';
import { InvoiceStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';

import previewStamp from '../assets/preview-stamp.png';

type InvoiceHeaderProps = {
  invoiceNumber: string;
  status?: InvoiceStatus | null;
};

export const InvoiceHeader: FC<InvoiceHeaderProps> = observer(
  ({ invoiceNumber, status }) => {
    const store = useStore();

    const isPreview = status === InvoiceStatus.Scheduled || !status;

    return (
      <div>
        <div className='flex flex-1 justify-between items-center'>
          <div className='flex items-center'>
            <h1 className='text-3xl font-bold'>Invoice</h1>
            {isPreview && (
              <img
                src={previewStamp}
                width={95}
                height={35}
                alt='Preview Stamp'
                className='absolute left-[6.5rem] top-2 rotate-[-10deg]'
              />
            )}
            {status && !isPreview && (
              <div className='ml-4 mt-1'>
                <Tag variant='outline' colorScheme='gray'>
                  {status}
                </Tag>
              </div>
            )}
          </div>

          {store.settings.tenant.value?.logoRepositoryFileId && (
            <div className='flex relative max-h-[120px] w-full justify-end'>
              <Image
                src={store.settings.tenant.value?.logoRepositoryFileId}
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

        {!isPreview && (
          <h2 className='text-sm text-gray-500'>N° {invoiceNumber}</h2>
        )}
      </div>
    );
  },
);
