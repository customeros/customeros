import { useNavigate, useSearchParams } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button';
import { Skeleton } from '@ui/feedback/Skeleton';
import { Download02 } from '@ui/media/icons/Download02';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
import { InvoicePreviewModalContent } from '@shared/components/Invoice/InvoicePreviewModal';
import {
  Modal,
  ModalPortal,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';

import { PaymentStatusSelect } from '../shared';

export const Preview = () => {
  const navigate = useNavigate();
  const client = getGraphQLClient();
  const [searchParams] = useSearchParams();
  const invoiceId = searchParams?.get('preview');

  const onOpenChange = () => {
    const newParams = new URLSearchParams(searchParams?.toString());
    newParams.delete('preview');
    window.history.pushState({}, '', `?${newParams.toString()}`);
    navigate(`?${newParams.toString()}`);
  };

  const { data, isLoading, isError } = useGetInvoiceQuery(
    client,
    {
      id: invoiceId ?? '',
    },
    {
      enabled: !!invoiceId,
    },
  );

  return (
    <Modal open={!!invoiceId} onOpenChange={onOpenChange}>
      <ModalPortal>
        <ModalOverlay>
          {/* width and height of A4 */}
          <ModalContent className='max-w-[794px]'>
            <ModalHeader className='flex justify-between items-center py-3 px-4'>
              {data?.invoice?.status ? (
                <PaymentStatusSelect
                  variant='invoice-preview'
                  value={data?.invoice?.status}
                  invoiceId={data?.invoice?.metadata?.id}
                />
              ) : (
                <Skeleton className='w-[72px] h-[34px]' />
              )}

              <Button
                asChild
                className='rounded-full'
                leftIcon={<Download02 />}
              >
                <a
                  target='_blank'
                  rel='noopener noreferrer'
                  href={`/fs/file/${invoiceId}/download?inline=true`}
                  className='no-underline text-gray-700 hover:no-underline'
                >
                  Download
                </a>
              </Button>
            </ModalHeader>
            <div className='h-[1123px]'>
              <InvoicePreviewModalContent
                data={data}
                isError={isError}
                isFetching={isLoading}
              />
            </div>
          </ModalContent>
        </ModalOverlay>
      </ModalPortal>
    </Modal>
  );
};
