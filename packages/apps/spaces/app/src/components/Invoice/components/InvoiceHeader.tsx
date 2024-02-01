'use client';

import React, { FC } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Image } from '@ui/media/Image';
import { Tag } from '@ui/presentation/Tag';
import { Heading } from '@ui/typography/Heading';

type InvoiceHeaderProps = {
  status?: string;
  invoiceNumber: string;
  logoUrl?: string | null;
};

export const InvoiceHeader: FC<InvoiceHeaderProps> = ({
  invoiceNumber,
  status,
  logoUrl,
}) => (
  <Box>
    <Flex alignItems='center' flex={1} justifyContent='space-between'>
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

      {logoUrl && (
        <Box position='relative' width='136px' height='30px'>
          <Image
            alt='' // Leaving alt empty cause provider info is available in the invoice and logo serves as a decoration
            src={logoUrl}
            fill
            style={{ objectFit: 'contain' }}
          />
        </Box>
      )}
    </Flex>

    <Heading as='h2' fontSize='sm' fontWeight='regular' color='gray.500'>
      N° {invoiceNumber}
    </Heading>
  </Box>
);
