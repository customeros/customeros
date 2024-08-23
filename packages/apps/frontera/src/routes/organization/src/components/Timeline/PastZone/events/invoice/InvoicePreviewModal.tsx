import { observer } from 'mobx-react-lite';

import { XClose } from '@ui/media/icons/XClose';
import { useStore } from '@shared/hooks/useStore';
import { Link01 } from '@ui/media/icons/Link01.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { InvoiceWithId } from '@organization/components/Timeline/types';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { InvoiceActionHeader } from '@shared/components/Invoice/InvoiceActionHeader';
import { InvoicePreviewModalContent } from '@shared/components/Invoice/InvoicePreviewModal';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const InvoicePreviewModal = observer(() => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const [_, copy] = useCopyToClipboard();
  const event = modalContent as Pick<InvoiceWithId, 'id' | '__typename'>;

  const store = useStore();
  const invoice = store.invoices.value.get(event.id);

  return (
    <>
      <CardHeader
        onClick={(e) => e.stopPropagation()}
        className='py-4 px-6 pb-1 bg-white top-0 rounded-xl'
      >
        <div className='flex justify-between items-center'>
          <InvoiceActionHeader
            status={invoice?.value?.status}
            id={invoice?.value?.metadata?.id}
            number={invoice?.value?.invoiceNumber}
          />

          <div className='flex justify-end items-center'>
            <Tooltip side='bottom' asChild={false} label='Copy invoice link'>
              <IconButton
                size='xs'
                variant='ghost'
                className='mr-1'
                colorScheme='gray'
                aria-label='Copy invoice link'
                icon={<Link01 color='text-inherit' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip
              label='Close'
              side='bottom'
              asChild={false}
              aria-label='close'
            >
              <IconButton
                size='xs'
                variant='ghost'
                colorScheme='gray'
                onClick={closeModal}
                aria-label='Close preview'
                icon={<XClose color='text-inherit' />}
              />
            </Tooltip>
          </div>
        </div>
      </CardHeader>

      <Card className='flex flex-col m-6 mt-3 p-4 shadow-xs h-full max-h-[80vh] overflow-y-auto'>
        <CardContent className='flex flex-1 p-0 items-center'>
          <InvoicePreviewModalContent
            invoiceStore={invoice}
            isFetching={store.invoices.isLoading}
          />
        </CardContent>
      </Card>
    </>
  );
});
