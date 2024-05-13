import { InvoiceStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { Download02 } from '@ui/media/icons/Download02';
import { renderStatusNode } from '@shared/components/Invoice/Cells';
import { StatusMenuButton } from '@shared/components/Invoice/components/StatusMenuButton';

import { DownloadInvoice } from '../../../../services/fileStore';

type InvoiceProps = {
  id?: string | null;
  number?: string | null;
  status?: InvoiceStatus | null;
};

export function InvoiceActionHeader({ status, id, number }: InvoiceProps) {
  const handleDownload = () => {
    if (!id || !number) {
      throw Error('Invoice cannot be downloaded without id or number');
    }

    return DownloadInvoice(id, number);
  };

  return (
    <div className='flex justify-between w-full'>
      {id ? (
        <StatusMenuButton status={status} id={id} />
      ) : (
        <div className='flex items-center'>{renderStatusNode(status)}</div>
      )}

      <div className='flex'>
        <Button
          variant='outline'
          colorScheme='gray'
          size='xs'
          className='rounded-full mr-2 bg-white'
          leftIcon={<Download02 className='size-3' />}
          onClick={handleDownload}
        >
          Download
        </Button>
      </div>
    </div>
  );
}
