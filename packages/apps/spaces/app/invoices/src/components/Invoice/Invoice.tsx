'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Tag } from '@ui/presentation/Tag';
import { Heading } from '@ui/typography/Heading';
import { Divider } from '@ui/presentation/Divider';
import { FileDownload02 } from '@ui/media/icons/FileDownload02';
import { ReceiptDownload } from '@ui/media/icons/ReceiptDownload';

import { ServicesTable } from './ServicesTable';
type Address = {
  suite: string;
  email: string;
  street: string;
};

type Service = {
  name: string;
  quantity: number;
  unitPrice: number;
};

type InvoiceProps = {
  tax: number;
  note: string;
  from: Address;
  total: number;
  dueDate: string;
  subtotal: number;
  lines: Service[];
  issueDate: string;
  billedTo: Address;
  invoiceNumber: string;
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
}: InvoiceProps) {
  return (
    <Flex width='80%' px={4} flexDir='column'>
      <Flex justifyContent='space-between' py={3}>
        <Tag
          colorScheme='success'
          variant='outline'
          borderRadius='md'
          boxShadow='unset'
          border='1px solid'
          borderColor='gray.300'
          cursor='pointer'
        >
          <Text>Paid</Text>
        </Tag>
        <Flex>
          <Button
            variant='outline'
            size='sm'
            borderRadius='full'
            leftIcon={<FileDownload02 />}
            mr={2}
          >
            Invoice
          </Button>
          <Button
            variant='outline'
            size='sm'
            borderRadius='full'
            leftIcon={<ReceiptDownload />}
          >
            Receipt
          </Button>
        </Flex>
      </Flex>

      <Flex flexDir='column' mt={2}>
        <Heading as='h1' fontSize='3xl' fontWeight='bold'>
          Invoice
        </Heading>
        <Heading as='h2' fontSize='md' color='gray.500'>
          N° {invoiceNumber}
        </Heading>

        <Flex
          mt={2}
          borderTop='1px solid'
          borderBottom='1px solid'
          borderColor='gray.300'
          justifyContent='space-evenly'
        >
          <Flex flexDir='column' flex={1} minW={150} py={2}>
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Issued
            </Text>
            <Text fontSize='sm' mb={4}>
              {issueDate}
            </Text>
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Due
            </Text>
            <Text fontSize='sm'>{dueDate}</Text>
          </Flex>
          <Divider orientation='vertical' mr={3} />
          <Flex flexDir='column' flex={1} minW={150} py={2}>
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Billed to
            </Text>
            <Text fontSize='sm' fontWeight='medium' mb={1}>
              {billedTo.street}
            </Text>
            <Text fontSize='sm'>{billedTo.suite}</Text>
            <Text fontSize='sm'>{billedTo.email}</Text>
          </Flex>
          <Divider orientation='vertical' mr={3} />
          <Flex flexDir='column' flex={1} minW={150} py={2}>
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              From
            </Text>
            <Text fontSize='sm' fontWeight='medium' mb={1}>
              {from.street}
            </Text>
            <Text fontSize='sm'>{from.suite}</Text>
            <Text fontSize='sm'>{from.email}</Text>
          </Flex>
        </Flex>
      </Flex>

      <Flex mt={4} flexDir='column'>
        <ServicesTable services={lines} />
        <Flex flexDir='column' alignSelf='flex-end' w='50%' maxW='300px' mt={4}>
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='medium'>
              Subtotal
            </Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              ${subtotal.toFixed(2)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.300' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm'>Tax</Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              ${tax.toFixed(2)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.300' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='medium'>
              Total
            </Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              ${total.toFixed(2)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.500' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='semibold'>
              Amount due
            </Text>
            <Text fontSize='sm' fontWeight='semibold' ml={2}>
              ${total.toFixed(2)}
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.500' />
          <Flex>
            <Text fontSize='sm' fontWeight='medium'>
              Note:
            </Text>
            <Text fontSize='sm' ml={2} color='gray.500'>
              {note}
            </Text>
          </Flex>
        </Flex>
      </Flex>
    </Flex>
  );
}
