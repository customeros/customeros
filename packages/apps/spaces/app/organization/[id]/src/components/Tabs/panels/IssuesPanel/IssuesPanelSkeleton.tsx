'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Card } from '@ui/layout/Card';
import { CardBody } from '@ui/presentation/Card';
import { Skeleton, SkeletonCircle } from '@ui/presentation/Skeleton';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const IssuesPanelSkeleton = () => {
  return (
    <OrganizationPanel title='Issues'>
      <Flex w='full' justify='flex-start'>
        <Skeleton
          borderRadius='full'
          w='50px'
          mb='1'
          h='16px'
          startColor='gray.300'
          endColor='gray.100'
        />
      </Flex>
      {Array.from({ length: 3 }).map((_, i) => (
        <Card
          key={i}
          w='full'
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
          <CardBody>
            <Flex flex='1' gap='4' alignItems='flex-start' flexWrap='wrap'>
              <SkeletonCircle
                height={10}
                width={10}
                startColor='gray.300'
                endColor='gray.100'
              />

              <Flex direction='column' flex={1}>
                <Flex justifyContent='space-between'>
                  <Skeleton
                    borderRadius='full'
                    h={3}
                    w={'50%'}
                    mb={2}
                    startColor='gray.300'
                    endColor='gray.100'
                  />
                </Flex>

                <Skeleton
                  borderRadius='full'
                  h={3}
                  w={'55%'}
                  mb={2}
                  startColor='gray.300'
                  endColor='gray.100'
                />
                <Skeleton
                  borderRadius='full'
                  h={3}
                  w={'45%'}
                  startColor='gray.300'
                  endColor='gray.100'
                />
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
          </CardBody>
        </Card>
      ))}
    </OrganizationPanel>
  );
};
