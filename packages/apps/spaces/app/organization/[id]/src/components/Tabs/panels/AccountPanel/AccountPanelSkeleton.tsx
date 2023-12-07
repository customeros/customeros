import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Divider } from '@ui/presentation/Divider';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { Skeleton, SkeletonCircle } from '@ui/presentation/Skeleton';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const AccountPanelSkeleton: React.FC = () => {
  return (
    <OrganizationPanel title='Account'>
      <Flex justify='space-between' w='full' align='center' mb={4}>
        <Flex justify='space-between' w='full' align='center'>
          <SkeletonCircle size='10' startColor='gray.300' endColor='gray.100' />
          <Flex
            ml='5'
            flexDir='column'
            align='flex-start'
            gap='1'
            flex='1'
            w='full'
          >
            <Skeleton
              w='45%'
              h='4'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
          <Skeleton
            w='55px'
            h='4'
            borderRadius='full'
            startColor='gray.300'
            endColor='gray.100'
          />
        </Flex>
      </Flex>

      <SkeletonCard>
        <SkeletonCardFooter1 />
      </SkeletonCard>
      <SkeletonCard>
        <SkeletonCardFooter2 />
      </SkeletonCard>
    </OrganizationPanel>
  );
};

const SkeletonCard = ({ children }: { children?: React.ReactNode }) => {
  return (
    <Card
      size='sm'
      width='full'
      borderRadius='xl'
      border='1px solid'
      borderColor='gray.200'
      boxShadow='xs'
      p='0'
      mb={4}
    >
      <CardBody as={Flex} align='center' w='full' p='4'>
        <Flex justify='space-between' w='full' align='center'>
          <SkeletonCircle size='10' startColor='gray.300' endColor='gray.100' />
          <Flex
            ml='5'
            flexDir='column'
            align='flex-start'
            gap='1'
            flex='1'
            w='full'
          >
            <Skeleton
              w='45%'
              h='4'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              w='35%'
              h='3'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
        </Flex>
      </CardBody>

      {children}
    </Card>
  );
};

const SkeletonCardFooter1 = () => {
  return (
    <CardFooter as={Flex} flexDir='column' p='4' pt='0'>
      <Flex justify='space-between' gap='4' align='center' w='full'>
        <VStack spacing='1' flex='1' align='flex-start'>
          <Skeleton
            w='65%'
            h='3'
            borderRadius='full'
            startColor='gray.300'
            endColor='gray.100'
          />
          <Flex w='full' gap='3' align='center' h='10'>
            <SkeletonCircle
              h='5'
              w='5'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              w='full'
              h='4'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
        </VStack>

        <VStack spacing='1' flex='1' align='flex-start'>
          <Skeleton
            w='65%'
            h='3'
            borderRadius='full'
            startColor='gray.300'
            endColor='gray.100'
          />
          <Flex w='full' gap='3' align='center' h='10'>
            <SkeletonCircle
              w='5'
              h='5'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              w='full'
              h='4'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
        </VStack>
      </Flex>
      <Flex justify='space-between' gap='4' align='center' w='full'>
        <VStack spacing='1' flex='1' align='flex-start'>
          <Skeleton
            w='65%'
            h='3'
            borderRadius='full'
            startColor='gray.300'
            endColor='gray.100'
          />
          <Flex w='full' gap='3' align='center' h='10'>
            <SkeletonCircle
              h='5'
              w='5'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              w='full'
              h='4'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
        </VStack>

        <VStack spacing='1' flex='1' align='flex-start'>
          <Skeleton
            w='65%'
            h='3'
            borderRadius='full'
            startColor='gray.300'
            endColor='gray.100'
          />
          <Flex w='full' gap='3' align='center' h='10'>
            <SkeletonCircle
              w='5'
              h='5'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              w='full'
              h='4'
              borderRadius='full'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
        </VStack>
      </Flex>
    </CardFooter>
  );
};

const SkeletonCardFooter2 = () => {
  return (
    <CardFooter as={Flex} flexDir='column' p='4' pt='0'>
      <Divider mb='4' mt='0' />

      <Flex w='full' gap='1' align='center'>
        <SkeletonCircle w='5' h='5' startColor='gray.300' endColor='gray.100' />
        <Skeleton
          w='45%'
          h='3'
          borderRadius='full'
          startColor='gray.300'
          endColor='gray.100'
        />
      </Flex>
    </CardFooter>
  );
};
