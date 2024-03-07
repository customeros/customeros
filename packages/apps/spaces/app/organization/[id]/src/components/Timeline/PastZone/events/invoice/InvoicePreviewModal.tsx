import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/presentation/Tooltip';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { useGetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
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
        py='4'
        px='6'
        pb='1'
        position='sticky'
        background='white'
        top={0}
        borderRadius='xl'
        onClick={(e) => e.stopPropagation()}
      >
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <InvoiceActionHeader
            status={data?.invoice?.status}
            id={data?.invoice?.metadata?.id}
            number={data?.invoice?.invoiceNumber}
          />

          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy invoice link' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy invoice link'
                color='gray.500'
                size='sm'
                mr={1}
                icon={<Link03 color='gray.500' height='18px' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                size='sm'
                icon={<XClose color='gray.500' height='24px' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>

      <Card m={6} mt={3} p='4' boxShadow='xs' variant='outline' w={600}>
        <CardBody as={Flex} p='0' align='center'>
          <InvoicePreviewModalContent
            data={data}
            isFetching={isFetching}
            isError={isError}
          />
        </CardBody>
      </Card>
    </>
  );
};
