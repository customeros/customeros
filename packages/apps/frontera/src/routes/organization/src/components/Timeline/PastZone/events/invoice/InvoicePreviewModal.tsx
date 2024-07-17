import { observer } from 'mobx-react-lite';

import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { useStore } from '@shared/hooks/useStore';
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
        className='py-4 px-6 pb-1 bg-white top-0 rounded-xl'
        onClick={(e) => e.stopPropagation()}
      >
        <div className='flex justify-between items-center'>
          <InvoiceActionHeader
            status={invoice?.value?.status}
            id={invoice?.value?.metadata?.id}
            number={invoice?.value?.invoiceNumber}
          />

          <div className='flex justify-end items-center'>
            <Tooltip label='Copy invoice link' side='bottom' asChild={false}>
              <IconButton
                className='mr-1'
                variant='ghost'
                aria-label='Copy invoice link'
                colorScheme='gray'
                size='md'
                icon={<Link03 color='gray.500' height='18px' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip
              label='Close'
              aria-label='close'
              side='bottom'
              asChild={false}
            >
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                colorScheme='gray'
                size='md'
                icon={<XClose color='gray.500' height='24px' />}
                onClick={closeModal}
              />
            </Tooltip>
          </div>
        </div>
      </CardHeader>

      <Card className='flex flex-col m-6 mt-3 p-4 shadow-xs w-[600px] h-[100%] overflow-y-auto'>
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
