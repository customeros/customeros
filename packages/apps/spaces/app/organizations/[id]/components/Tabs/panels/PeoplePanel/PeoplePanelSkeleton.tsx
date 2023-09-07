'use client';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Icons } from '@ui/media/Icon';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { CardHeader } from '@ui/presentation/Card';
import { Card } from '@ui/layout/Card';
import { Skeleton, SkeletonCircle } from '@ui/presentation/Skeleton';

export const PeoplePanelSkeleton = () => {
  return (
    <OrganizationPanel
      title='People'
      actionItem={
        <Button
          size='sm'
          variant='outline'
          leftIcon={<Icons.UsersPlus color='gray.500' />}
          type='button'
          isDisabled
        >
          Add
        </Button>
      }
    >
      {Array.from({ length: 3 }).map((_, i) => (
        <Card
          key={i}
          w='full'
          minH='106px'
          boxShadow={'xs'}
          cursor='pointer'
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
          <CardHeader as={Flex} p='4' pb={2} position='relative'>
            <SkeletonCircle
              height={12}
              width={12}
              boxShadow='avatarRingGray'
              startColor='gray.300'
              endColor='gray.100'
            />

            <Flex ml='4' flexDir='column' flex='1'>
              <Skeleton
                borderRadius='full'
                h={3}
                w={100}
                mb={3}
                startColor='gray.300'
                endColor='gray.100'
              />
              <Skeleton
                borderRadius='full'
                h={3}
                w={200}
                mb={4}
                startColor='gray.300'
                endColor='gray.100'
              />
              <Skeleton
                borderRadius='full'
                h={3}
                w={250}
                mb={2}
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
