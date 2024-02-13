'use client';
import React from 'react';

import { useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Link03 } from '@ui/media/icons/Link03';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/presentation/Tooltip';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { useGetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
import { InvoiceActionHeader } from '@shared/components/Invoice/InvoiceActionHeader';
import { InvoicePreviewModalContent } from '@shared/components/Invoice/InvoicePreviewModal';

export const PreviewPanel = ({ id }: { id: string }) => {
  const [_, copy] = useCopyToClipboard();
  const isRestoring = useIsRestoring();
  const client = getGraphQLClient();
  const { data, isFetching, isError } = useGetInvoiceQuery(client, {
    id,
  });

  return (
    <Flex flexDir='column'>
      <Flex
        direction='row'
        justifyContent='space-between'
        alignItems='center'
        px={4}
      >
        <InvoiceActionHeader
          status={data?.invoice?.status}
          id={data?.invoice?.id}
          number={data?.invoice?.number}
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
        </Flex>
      </Flex>

      <InvoicePreviewModalContent
        data={data}
        isFetching={isFetching || isRestoring}
        isError={isError}
      />
    </Flex>
  );
};
