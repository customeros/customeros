'use client';

import React, { FC } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Tag } from '@ui/presentation/Tag';
import { Heading } from '@ui/typography/Heading';
import { Image as ChakraImage } from '@ui/media/Image';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

type InvoiceHeaderProps = {
  status?: string;
  invoiceNumber: string;
};

export const InvoiceHeader: FC<InvoiceHeaderProps> = ({
  invoiceNumber,
  status,
}) => {
  const client = getGraphQLClient();

  const { data: globalCacheData } = useGlobalCacheQuery(client);

  return (
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

        {globalCacheData?.global_Cache?.cdnLogoUrl && (
          <Flex
            position='relative'
            maxHeight={120}
            width='full'
            justifyContent='flex-end'
          >
            <ChakraImage
              src={`${globalCacheData?.global_Cache?.cdnLogoUrl}`}
              alt='CustomerOS'
              width={136}
              height={40}
              style={{
                objectFit: 'contain',
                maxHeight: '40px',
                maxWidth: 'fit-content',
              }}
            />
          </Flex>
        )}
      </Flex>

      <Heading as='h2' fontSize='sm' fontWeight='regular' color='gray.500'>
        N° {invoiceNumber}
      </Heading>
    </Box>
  );
};
