import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { Download02 } from '@ui/media/icons/Download02';
import { DownloadFile } from '@ui/media/DownloadFile/DownloadFile';
import { InvoicePreviewModalContent } from '@shared/components/Invoice/InvoicePreviewModal';
import {
  Modal,
  ModalPortal,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';

import { PaymentStatusSelect } from '../shared';

export const Preview = observer(() => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const invoiceNumber = searchParams?.get('preview');

  const onOpenChange = () => {
    const newParams = new URLSearchParams(searchParams?.toString());

    newParams.delete('preview');
    window.history.pushState({}, '', `?${newParams.toString()}`);
    navigate(`?${newParams.toString()}`);
  };

  const store = useStore();

  const invoice = invoiceNumber
    ? store.invoices.value.get(invoiceNumber)
    : null;

  return (
    <Modal open={!!invoiceNumber} onOpenChange={onOpenChange}>
      <ModalPortal>
        <ModalOverlay className='z-50'>
          {/* width and height of A4 */}
          <ModalContent className='max-w-[794px]'>
            <ModalHeader className='flex justify-between items-center py-3 px-4'>
              {invoice?.value?.status ? (
                <PaymentStatusSelect
                  variant='invoice-preview'
                  invoiceNumber={invoice?.number}
                />
              ) : (
                <Skeleton className='w-[72px] h-[34px]' />
              )}

              <DownloadFile
                variant='outline'
                leftIcon={<Download02 />}
                fileId={invoice?.id ?? ''}
                fileName={`invoice-${invoice?.value?.invoiceNumber}`}
              />
            </ModalHeader>
            <div className='h-[1123px]'>
              <InvoicePreviewModalContent
                invoiceStore={invoice}
                isFetching={store.invoices.isLoading}
              />
            </div>
          </ModalContent>
        </ModalOverlay>
      </ModalPortal>
    </Modal>
  );
});
