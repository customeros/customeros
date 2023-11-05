'use client';
import Image from 'next/image';
import React, { FC } from 'react';

import { Box } from '@ui/layout/Box';
import { VStack } from '@ui/layout/Stack';
import { Skeleton } from '@ui/presentation/Skeleton';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';

export const TimelineItemSkeleton: FC = () => {
  return (
    <Box mt={4} mr={6}>
      <Skeleton
        height='0.5rem'
        width='100px'
        borderRadius='md'
        mb={4}
        startColor='gray.300'
        endColor='gray.100'
      />
      <Card
        variant='outline'
        size='md'
        fontSize='14px'
        background='white'
        flexDirection='row'
        maxWidth={549}
        position='unset'
        aspectRatio='9/2'
        boxShadow='xs'
      >
        <CardBody
          pt={5}
          pb={5}
          pl={5}
          pr={0}
          overflow={'hidden'}
          flexDirection='row'
        >
          <VStack
            align='flex-start'
            spacing={0}
            justifyContent='space-between'
            h='100%'
          >
            <Skeleton
              width='33%'
              height='0.75rem'
              borderRadius='md'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              width='95%'
              height='0.5rem'
              borderRadius='md'
              startColor='gray.300'
              endColor='gray.100'
              // mb={1}
            />

            <Skeleton
              width='95%'
              height='0.5rem'
              borderRadius='md'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              width='95%'
              height='0.5rem'
              borderRadius='md'
              startColor='gray.300'
              endColor='gray.100'
            />
          </VStack>
        </CardBody>
        <CardFooter pt={5} pb={5} pr={5} pl={0} ml={1}>
          <Skeleton
            h='70px'
            w='54px'
            borderRadius='4px'
            startColor='gray.300'
            endColor='gray.100'
          />
        </CardFooter>
      </Card>
    </Box>
  );
};
