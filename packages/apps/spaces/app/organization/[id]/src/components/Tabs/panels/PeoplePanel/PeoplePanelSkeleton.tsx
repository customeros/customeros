'use client';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { Skeleton } from '@ui/feedback/Skeleton/Skeleton';
import { Card, CardHeader } from '@ui/presentation/Card/Card';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

export const PeoplePanelSkeleton = () => {
  return (
    <OrganizationPanel
      title='People'
      actionItem={
        <Button
          size='sm'
          variant='outline'
          leftIcon={<Icons.UsersPlus color='gray.500' />}
          isDisabled
        >
          Add
        </Button>
      }
    >
      {Array.from({ length: 3 }).map((_, i) => (
        <Card
          className='bg-white w-full min-h-[106px] group cursor-pointer border-1 border-gray-200 rounded-lg shadow-xs'
          key={i}
        >
          <CardHeader className='flex p-4 pb-2 relative'>
            <Skeleton className='size-12 rounded-full ring-offset-1 ring-4 ring-offset-gray-200 ring-gray-100' />

            <Flex ml='4' flexDir='column' flex='1'>
              <Skeleton className='h-3 w-[100px] mb-3 rounded-full' />
              <Skeleton className='h-3 w-[200px] mb-4 rounded-full' />
              <Skeleton className='h-3 w-[250px] mb-3 rounded-full' />
            </Flex>
          </CardHeader>
        </Card>
      ))}
    </OrganizationPanel>
  );
};
