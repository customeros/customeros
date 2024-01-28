'use client';

import React from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Tag } from '@ui/presentation/Tag';
import { Text } from '@ui/typography/Text';
import { InvoiceLine } from '@graphql/types';
import { Heading } from '@ui/typography/Heading';
import { DateTimeUtils } from '@spaces/utils/date';
import { Divider } from '@ui/presentation/Divider';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { ServicesTable } from './ServicesTable';
type Address = {
  zip: string;
  email: string;
  name?: string;
  country: string;
  locality: string;
  addressLine: string;
  addressLine2?: string;
};

type InvoiceProps = {
  tax: number;
  note: string;
  from: Address;
  total: number;
  dueDate: string;
  status?: string;
  subtotal: number;
  currency?: string;
  issueDate: string;
  billedTo: Address;
  amountDue?: number;
  invoiceNumber: string;
  lines: Array<InvoiceLine>;
  isBilledToFocused?: boolean;
};

export function Invoice({
  invoiceNumber,
  issueDate,
  dueDate,
  billedTo,
  from,
  lines,
  subtotal,
  tax,
  total,
  note,
  amountDue,
  status,
  isBilledToFocused,
  currency = 'USD',
}: InvoiceProps) {
  return (
    <Flex px={4} flexDir='column' w='inherit' overflowY='auto'>
      <Flex flexDir='column' mt={2}>
        <Flex alignItems='center'>
          <Heading as='h1' fontSize='3xl' fontWeight='bold'>
            Invoice
          </Heading>
          {status && (
            <Box ml={4}>
              <Tag variant='outline' colorScheme='gray'>
                {status}
              </Tag>
            </Box>
          )}
        </Flex>

        <Heading as='h2' fontSize='sm' fontWeight='regular' color='gray.500'>
          N° {invoiceNumber}
        </Heading>

        <Flex
          mt={2}
          borderTop='1px solid'
          borderBottom='1px solid'
          borderColor='gray.300'
          justifyContent='space-evenly'
          gap={3}
        >
          <Flex
            flexDir='column'
            flex={1}
            minW={150}
            py={2}
            borderRight={'1px solid'}
            filter={isBilledToFocused ? 'blur(2px)' : 'none'}
            borderColor='gray.300'
          >
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Issued
            </Text>
            <Text fontSize='sm' mb={4}>
              {DateTimeUtils.format(
                issueDate,
                DateTimeUtils.dateWithAbreviatedMonth,
              )}
            </Text>
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Due
            </Text>
            <Text fontSize='sm'>
              {DateTimeUtils.format(
                dueDate,
                DateTimeUtils.dateWithAbreviatedMonth,
              )}
            </Text>
          </Flex>
          <Flex
            flexDir='column'
            minW={150}
            w='160px'
            py={2}
            pr={2}
            borderRight={isBilledToFocused ? '1px solid' : '1px solid'}
            borderColor={'gray.300'}
            position='relative'
            sx={{
              '&:after': {
                content: '""',
                bg: 'transparent',
                border: '2px solid',
                position: 'absolute',
                top: 0,
                bottom: 0,
                left: -3,
                right: 0,
                opacity: isBilledToFocused ? 1 : 0,
              },
            }}
          >
            <Text fontWeight='semibold' mb={0.5} fontSize='sm'>
              Billed to
            </Text>
            <Text fontSize='sm' fontWeight='medium' mb={1} lineHeight={1.2}>
              {billedTo.name}
            </Text>

            <Text fontSize='sm' lineHeight={1.2}>
              {billedTo.addressLine}
              <Text as='span' display='block' lineHeight={1.2}>
                {billedTo.addressLine2}
              </Text>
            </Text>
            <Text fontSize='sm' lineHeight={1.2}>
              {billedTo.locality} {billedTo.locality && ', '} {billedTo.zip}
            </Text>
            <Text fontSize='sm' lineHeight={1.2}>
              {billedTo.country}
            </Text>
            <Text fontSize='sm' lineHeight={1.2}>
              {billedTo.email}
            </Text>
          </Flex>
          <Flex
            flexDir='column'
            flex={1}
            minW={150}
            py={2}
            filter={isBilledToFocused ? 'blur(2px)' : 'none'}
          >
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              From
            </Text>
            <Text fontSize='sm' fontWeight='medium' mb={1} lineHeight={1.2}>
              {from.name}
            </Text>

            <Text fontSize='sm' lineHeight={1.2}>
              {from.addressLine}
              <Text as='span' display='block' lineHeight={1.2}>
                {from.addressLine2}
              </Text>
            </Text>
            <Text fontSize='sm' lineHeight={1.2}>
              {from.locality} {from.locality && ', '} {from.zip}
            </Text>
            <Text fontSize='sm' lineHeight={1.2}>
              {from.country}
            </Text>
            <Text fontSize='sm' lineHeight={1.2}>
              {from.email}
            </Text>
          </Flex>
        </Flex>
      </Flex>

      <Flex
        mt={4}
        flexDir='column'
        filter={isBilledToFocused ? 'blur(2px)' : 'none'}
      >
        <ServicesTable services={lines} currency={currency} />
        <Flex flexDir='column' alignSelf='flex-end' w='50%' maxW='300px' mt={4}>
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='medium'>
              Subtotal
            </Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              {formatCurrency(subtotal, 2, currency)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.300' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm'>Tax</Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              {formatCurrency(tax, 2, currency)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.300' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='medium'>
              Total
            </Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              {formatCurrency(total, 2, currency)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.500' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='semibold'>
              Amount due
            </Text>
            <Text fontSize='sm' fontWeight='semibold' ml={2}>
              {formatCurrency(amountDue || total, 2, currency)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.500' />

          {note && (
            <Flex>
              <Text fontSize='sm' fontWeight='medium'>
                Note:
              </Text>
              <Text fontSize='sm' ml={2} color='gray.500'>
                {note}
              </Text>
            </Flex>
          )}
        </Flex>
      </Flex>
    </Flex>
  );
}
