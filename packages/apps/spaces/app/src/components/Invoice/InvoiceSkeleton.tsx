'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { Divider } from '@ui/presentation/Divider';
import { Skeleton } from '@ui/presentation/Skeleton';

import { ServicesTable } from './ServicesTable';

export function InvoiceSkeleton() {
  return (
    <Flex px={4} flexDir='column' w='inherit'>
      <Flex flexDir='column' mt={2}>
        <Flex alignItems='center'>
          <Heading as='h1' fontSize='3xl' fontWeight='bold'>
            Invoice
          </Heading>
        </Flex>

        <Flex
          fontSize='sm'
          fontWeight='regular'
          color='gray.500'
          alignItems='center'
        >
          N° <Skeleton width='60px' height='12px' ml={1} />
        </Flex>

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
            borderColor='gray.300'
          >
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Issued
            </Text>
            <Skeleton width='50px' height='12px' />
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Due
            </Text>
            <Skeleton width='50px' height='12px' />
          </Flex>
          <Flex
            flexDir='column'
            flex={1}
            minW={150}
            py={2}
            borderColor={'gray.300'}
            position='relative'
          >
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              Billed to
            </Text>
            <Skeleton width='150px' height='12px' mb={1} />
            <Skeleton width='120px' height='12px' mb={1} />
            <Skeleton width='50px' height='12px' mb={1} />

            <Flex>
              <Skeleton width='60px' height='12px' mr={2} mb={1} />
              <Skeleton width='50px' height='12px' mb={1} />
            </Flex>
            <Skeleton width='40px' height='12px' mb={1} />
            <Skeleton width='90px' height='12px' />
          </Flex>
          <Flex flexDir='column' flex={1} minW={150} py={2}>
            <Text fontWeight='semibold' mb={1} fontSize='sm'>
              From
            </Text>
            <Skeleton width='150px' height='12px' mb={1} />
            <Skeleton width='120px' height='12px' mb={1} />
            <Skeleton width='50px' height='12px' mb={1} />

            <Flex>
              <Skeleton width='60px' height='12px' mr={2} mb={1} />
              <Skeleton width='50px' height='12px' mb={1} />
            </Flex>
            <Skeleton width='40px' height='12px' mb={1} />
            <Skeleton width='90px' height='12px' />
          </Flex>
        </Flex>
      </Flex>

      <Flex mt={4} flexDir='column'>
        <ServicesTable services={[]} currency='USD' />
        <Flex my={2} justifyContent='space-between'>
          <Skeleton width='55%' height='14px' mr={2} />
          <Skeleton width='10%' height='14px' mr={2} />
          <Skeleton width='20%' height='14px' mr={2} />
          <Skeleton width='15%' height='14px' mr={2} />
        </Flex>
        <Flex my={2} justifyContent='space-between'>
          <Skeleton width='55%' height='14px' mr={2} />
          <Skeleton width='10%' height='14px' mr={2} />
          <Skeleton width='20%' height='14px' mr={2} />
          <Skeleton width='15%' height='14px' mr={2} />
        </Flex>
        <Flex flexDir='column' alignSelf='flex-end' w='50%' maxW='300px' mt={4}>
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='medium'>
              Subtotal
            </Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              <Skeleton width='20px' height='12px' />
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.300' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm'>Tax</Text>
            <Text fontSize='sm' ml={2} color='gray.600'>
              <Skeleton width='20px' height='12px' />
            </Text>
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.300' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='medium'>
              Total
            </Text>
            <Skeleton width='20px' height='12px' />
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.500' />
          <Flex justifyContent='space-between'>
            <Text fontSize='sm' fontWeight='semibold'>
              Amount due
            </Text>
            <Skeleton width='20px' height='12px' />
          </Flex>
          <Divider orientation='horizontal' my={1} borderColor='gray.500' />
        </Flex>
      </Flex>
    </Flex>
  );
}
