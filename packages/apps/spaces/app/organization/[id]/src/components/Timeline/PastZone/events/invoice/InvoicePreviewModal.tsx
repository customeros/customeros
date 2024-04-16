import React from 'react';

import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { useGetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { InvoiceWithId } from '@organization/src/components/Timeline/types';
import { InvoiceActionHeader } from '@shared/components/Invoice/InvoiceActionHeader';
import { InvoicePreviewModalContent } from '@shared/components/Invoice/InvoicePreviewModal';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

export const InvoicePreviewModal = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const [_, copy] = useCopyToClipboard();
  const client = getGraphQLClient();
  const event = modalContent as Pick<InvoiceWithId, 'id' | '__typename'>;
  const { data, isFetching, isError } = useGetInvoiceQuery(client, {
    id: event.id,
  });

  return (
    <>
      <CardHeader
        className='py-4 px-6 pb-1 bg-white top-0 rounded-xl'
        onClick={(e) => e.stopPropagation()}
      >
        <div className='flex justify-between items-center'>
          <InvoiceActionHeader
            status={data?.invoice?.status}
            id={data?.invoice?.metadata?.id}
            number={data?.invoice?.invoiceNumber}
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
            data={data}
            isFetching={isFetching}
            isError={isError}
          />
        </CardContent>
      </Card>
    </>
  );
};
