'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { InvoiceStatus } from '@graphql/types';
import { Download02 } from '@ui/media/icons/Download02';
import { StatusCell } from '@shared/components/Invoice/Cells';

type InvoiceProps = {
  onDownload: () => void;
  status?: InvoiceStatus | null;
};

export function InvoiceActionHeader({ status, onDownload }: InvoiceProps) {
  return (
    <Flex justifyContent='space-between' w='full'>
      <StatusCell status={status} />

      <Flex>
        <Button
          variant='outline'
          size='sm'
          borderRadius='full'
          leftIcon={<Download02 />}
          onClick={onDownload}
          mr={2}
        >
          Download
        </Button>
      </Flex>
    </Flex>
  );
}
