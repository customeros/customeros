'use client';
import React, { FC } from 'react';
import { Box } from '@ui/layout/Box';
import { Skeleton } from '@ui/presentation/Skeleton';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { VStack } from '@ui/layout/Stack';
import Image from 'next/image';

export const TimelineItemSkeleton: FC = () => {
  return (
    <Box mt={4} mr={6}>
      <Skeleton
        height='0.75rem'
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
              width='80%'
              height='0.75rem'
              borderRadius='md'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              width='80%'
              height='0.75rem'
              borderRadius='md'
              startColor='gray.300'
              endColor='gray.100'
              mb={1}
            />

            <Skeleton
              width='90%'
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
          <Image
            src={'/backgrounds/organization/post-stamp.webp'}
            alt='Email'
            width={54}
            height={70}
            style={{
              filter:
                'drop-shadow(0px 0.5px 1px #D8D8D8) brightness(1) grayscale(1)',
            }}
          />
        </CardFooter>
      </Card>
    </Box>
  );
};
