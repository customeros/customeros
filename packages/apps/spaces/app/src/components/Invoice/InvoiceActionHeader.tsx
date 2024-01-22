'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { InvoiceStatus } from '@graphql/types';
import { Download02 } from '@ui/media/icons/Download02';
import { StatusCell } from '@shared/components/Invoice/Cells';

type InvoiceProps = {
  status: InvoiceStatus;
};

export function InvoiceActionHeader({ status }: InvoiceProps) {
  return (
    <Flex px={4} flexDir='column' w='inherit'>
      <Flex justifyContent='space-between' py={3}>
        <StatusCell status={status} />

        <Flex>
          <Button
            variant='outline'
            size='sm'
            borderRadius='full'
            leftIcon={<Download02 />}
            mr={2}
          >
            Download
          </Button>
        </Flex>
      </Flex>
    </Flex>
  );
}
