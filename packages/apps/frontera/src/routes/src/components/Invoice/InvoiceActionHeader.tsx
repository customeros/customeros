import { observer } from 'mobx-react-lite';
import { PaymentStatusSelect } from '@invoices/components/shared';

import { InvoiceStatus } from '@graphql/types';
import { Download02 } from '@ui/media/icons/Download02';
import { renderStatusNode } from '@shared/components/Invoice/Cells';
import { DownloadFile } from '@ui/media/DownloadFile/DownloadFile.tsx';

type InvoiceProps = {
  id?: string | null;
  number?: string | null;
  status?: InvoiceStatus | null;
};

export const InvoiceActionHeader = observer(
  ({ status, id, number }: InvoiceProps) => {
    return (
      <div className='flex justify-between w-full'>
        {status ? (
          <PaymentStatusSelect
            value={status}
            variant='invoice-preview'
            invoiceNumber={number ?? ''}
          />
        ) : (
          <div className='flex items-center'>{renderStatusNode(status)}</div>
        )}

        <div className='flex mr-2'>
          {id && number && (
            <DownloadFile
              fileId={id}
              variant='outline'
              leftIcon={<Download02 />}
              fileName={`invoice-${number}`}
            />
          )}
        </div>
      </div>
    );
  },
);
