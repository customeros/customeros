import { observer } from 'mobx-react-lite';

import { InvoiceStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Download02 } from '@ui/media/icons/Download02';
import { renderStatusNode } from '@shared/components/Invoice/Cells';
import { StatusMenuButton } from '@shared/components/Invoice/components/StatusMenuButton';

type InvoiceProps = {
  id?: string | null;
  number?: string | null;
  status?: InvoiceStatus | null;
};

export const InvoiceActionHeader = observer(
  ({ status, id, number }: InvoiceProps) => {
    const store = useStore();

    const handleDownload = () => {
      if (!id || !number) {
        throw Error('Invoice cannot be downloaded without id or number');
      }

      return store.files.downloadAttachment(id, number);
    };

    return (
      <div className='flex justify-between w-full'>
        {id ? (
          <StatusMenuButton id={id} status={status} />
        ) : (
          <div className='flex items-center'>{renderStatusNode(status)}</div>
        )}

        <div className='flex'>
          <Button
            size='xs'
            variant='outline'
            colorScheme='gray'
            onClick={handleDownload}
            className='rounded-full mr-2 bg-white'
            leftIcon={<Download02 className='size-3' />}
          >
            Download
          </Button>
        </div>
      </div>
    );
  },
);
