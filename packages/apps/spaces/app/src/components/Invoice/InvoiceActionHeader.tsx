'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { InvoiceStatus } from '@graphql/types';
import { Download02 } from '@ui/media/icons/Download02';
import { StatusCell } from '@shared/components/Invoice/Cells';

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
    <Flex justifyContent='space-between' w='full'>
      <StatusCell status={status} />

      <Flex>
        <Button
          variant='outline'
          size='xs'
          borderRadius='full'
          leftIcon={<Download02 boxSize={3} />}
          onClick={handleDownload}
          mr={2}
          px={2}
        >
          Download
        </Button>
      </Flex>
    </Flex>
  );
}
