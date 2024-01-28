'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
// import { Button } from '@ui/form/Button';
import { InvoiceStatus } from '@graphql/types';
// import { Download02 } from '@ui/media/icons/Download02';
import { StatusCell } from '@shared/components/Invoice/Cells';

type InvoiceProps = {
  onDownload: () => void;
  status?: InvoiceStatus | null;
};

export function InvoiceActionHeader({ status }: InvoiceProps) {
  return (
    <Flex justifyContent='space-between' w='full'>
      <StatusCell status={status} />

      <Flex>
        {/*<Button*/}
        {/*  variant='outline'*/}
        {/*  size='xs'*/}
        {/*  borderRadius='full'*/}
        {/*  leftIcon={<Download02 boxSize={3} />}*/}
        {/*  onClick={onDownload}*/}
        {/*  mr={2}*/}
        {/*  px={2}*/}
        {/*>*/}
        {/*  Download*/}
        {/*</Button>*/}
      </Flex>
    </Flex>
  );
}
