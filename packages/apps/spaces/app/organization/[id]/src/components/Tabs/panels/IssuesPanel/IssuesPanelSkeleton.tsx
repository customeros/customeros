'use client';

import { Flex } from '@ui/layout/Flex';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { CardHeader } from '@ui/presentation/Card';
import { Card } from '@ui/layout/Card';
import { Skeleton, SkeletonCircle } from '@ui/presentation/Skeleton';
import React from 'react';

export const IssuesPanelSkeleton = () => {
  return (
    <OrganizationPanel title='Issues'>
      {Array.from({ length: 3 }).map((_, i) => (
        <Card
          key={i}
          w='full'
          h='66px'
          boxShadow={'xs'}
          cursor='pointer'
          size='sm'
          borderRadius='lg'
          border='1px solid'
          borderColor='gray.200'
          _hover={{
            boxShadow: 'md',
            '& > div > #confirm-button': {
              opacity: '1',
              pointerEvents: 'auto',
            },
          }}
          transition='all 0.2s ease-out'
        >
          <CardHeader>
            <Flex flex='1' gap='4' alignItems='flex-start' flexWrap='wrap'>
              <SkeletonCircle
                height={12}
                width={12}
                startColor='gray.300'
                endColor='gray.100'
              />

              <Flex direction='column' flex={1}>
                <Flex justifyContent='space-between'>
                  <Skeleton
                    borderRadius='full'
                    h={3}
                    w={200}
                    mb={2}
                    startColor='gray.300'
                    endColor='gray.100'
                  />
                </Flex>

                <Skeleton
                  borderRadius='full'
                  h={3}
                  w={220}
                  mb={2}
                  startColor='gray.300'
                  endColor='gray.100'
                />
                {/* TODO uncomment commented out code as soon as COS-464 is merged */}
                {/*<Skeleton*/}
                {/*  borderRadius='full'*/}
                {/*  h={3}*/}
                {/*  w={180}*/}
                {/*  startColor='gray.300'*/}
                {/*  endColor='gray.100'*/}
                {/*/>*/}
              </Flex>
              <Skeleton
                display='block'
                position='static'
                borderRadius='md'
                h={6}
                w={10}
                startColor='gray.300'
                endColor='gray.100'
              />
            </Flex>
          </CardHeader>
        </Card>
      ))}
    </OrganizationPanel>
  );
};
