'use client';

import React, { FC } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Divider } from '@ui/presentation/Divider';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

interface InvoiceSummaryProps {
  tax: number;
  total: number;
  subtotal: number;
  currency: string;
  amountDue?: number;
  note?: string | null;
}

export const InvoiceSummary: FC<InvoiceSummaryProps> = ({
  subtotal,
  tax,
  total,
  amountDue,
  currency,
  note,
}) => {
  return (
    <Flex flexDir='column' alignSelf='flex-end' w='50%' maxW='300px' mt={4}>
      <Flex justifyContent='space-between'>
        <Text fontSize='sm' fontWeight='medium' alignItems='center'>
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
  );
};
