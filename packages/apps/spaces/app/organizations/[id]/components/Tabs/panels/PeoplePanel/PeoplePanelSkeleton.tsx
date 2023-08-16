'use client';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Icons } from '@ui/media/Icon';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { CardHeader } from '@ui/presentation/Card';
import { Avatar } from '@ui/media/Avatar';
import User from '@spaces/atoms/icons/User';
import { Card } from '@chakra-ui/card';
import { Skeleton } from '@chakra-ui/react';

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
      <Card
        w='full'
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
          <Avatar
            name=''
            variant='skeleton'
            icon={
              <User color={'var(--chakra-colors-gray-400)'} height='1.8rem' />
            }
          />
          <Flex ml='4' flexDir='column' flex='1'>
            <Skeleton h={3} w={100} mb={3} />
            <Skeleton h={3} w={200} mb={4} />
            <Skeleton h={3} w={250} mb={2} />
          </Flex>
        </CardHeader>
      </Card>
      <Card
        w='full'
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
          <Avatar
            name=''
            variant='skeleton'
            icon={
              <User color={'var(--chakra-colors-gray-400)'} height='1.8rem' />
            }
          />
          <Flex ml='4' flexDir='column' flex='1'>
            <Skeleton h={3} w={100} mb={3} />
            <Skeleton h={3} w={200} mb={4} />
            <Skeleton h={3} w={250} mb={2} />
          </Flex>
        </CardHeader>
      </Card>
    </OrganizationPanel>
  );
};
