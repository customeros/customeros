'use client';

import React, { FC } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Tag } from '@ui/presentation/Tag';
import { Heading } from '@ui/typography/Heading';

type InvoiceHeaderProps = {
  status?: string;
  invoiceNumber: string;
};

export const InvoiceHeader: FC<InvoiceHeaderProps> = ({
  invoiceNumber,
  status,
}) => (
  <Box>
    <Flex alignItems='center'>
      <Heading as='h1' fontSize='3xl' fontWeight='bold'>
        Invoice
      </Heading>
      {status && (
        <Box ml={4} mt={1}>
          <Tag variant='outline' colorScheme='gray'>
            {status}
          </Tag>
        </Box>
      )}
    </Flex>

    <Heading as='h2' fontSize='sm' fontWeight='regular' color='gray.500'>
      N° {invoiceNumber}
    </Heading>
  </Box>
);
