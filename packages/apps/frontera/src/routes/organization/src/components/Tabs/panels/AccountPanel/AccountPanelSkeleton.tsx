import React from 'react';

import { Skeleton } from '@ui/feedback/Skeleton';
import { Divider } from '@ui/presentation/Divider/Divider';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

export const AccountPanelSkeleton: React.FC = () => {
  return (
    <OrganizationPanel title='Account'>
      <div className='flex justify-between w-full items-center mb-4'>
        <div className='flex justify-between w-full items-center'>
          <Skeleton className='size-10 rounded-full' />
          <div className='flex ml-5 flex-col items-start gap-1 flex-1 w-full'>
            <Skeleton className='w-[45%] h-4 rounded-full' />
          </div>
          <Skeleton className='w-[55px] h-4 rounded-full' />
        </div>
      </div>

      <SkeletonCard showBtnSection>
        <SkeletonCardFooter1 />
      </SkeletonCard>
      <SkeletonCard>
        <SkeletonCardFooter2 />
      </SkeletonCard>
    </OrganizationPanel>
  );
};

const SkeletonCard = ({
  children,
  showBtnSection,
}: {
  showBtnSection?: boolean;
  children?: React.ReactNode;
}) => {
  return (
    <div className='w-full rounded-xl border border-gray-200 shadow-xs p-0 mb-4'>
      <div className='flex items-center w-full p-4'>
        <div className='flex justify-between w-full items-start'>
          <Skeleton className='size-10 rounded-full' />
          <div className='flex ml-5 flex-col items-start gap-1 flex-1 w-full'>
            <Skeleton className='w-[55%] h-4 rounded-full' />
            <Skeleton className='w-[35%] h-3 rounded-full' />
          </div>
          {showBtnSection && (
            <div className='flex items-start'>
              <Skeleton className='size-4 mr-2 rounded-md' />
              <Skeleton className='w-10 h-4 rounded-md' />
            </div>
          )}
        </div>
      </div>

      {children}
    </div>
  );
};

const SkeletonCardFooter1 = () => {
  return (
    <div className='flex flex-col p-4 pt-0'>
      <div className='flex justify-between gap-4 items-center w-full'>
        <div className='flex flex-col space-y-1 flex-1 items-start'>
          <Skeleton className='w-[65%] h-3 rounded-full' />
          <div className='flex w-full gap-3 items-center h-10'>
            <Skeleton className='size-5 rounded-full' />
            <Skeleton className='w-full h-4 rounded-full' />
          </div>
        </div>

        <div className='flex flex-col space-y-1 flex-1 items-start'>
          <Skeleton className='w-[65%] h-3 rounded-full' />
          <div className='flex w-full gap-3 items-center h-10'>
            <Skeleton className='size-5 rounded-full' />
            <Skeleton className='w-full h-4 rounded-full' />
          </div>
        </div>
      </div>
      <div className='flex justify-between gap-4 items-center w-full'>
        <div className='flex flex-col space-y-1 flex-1 items-start'>
          <Skeleton className='w-[65%] h-3 rounded-full' />
          <div className='flex w-full gap-3 items-center h-10'>
            <Skeleton className='size-5 rounded-full' />
            <Skeleton className='w-full h-4 rounded-full' />
          </div>
        </div>

        <div className='flex flex-col space-y-1 flex-1 items-start'>
          <Skeleton className='w-[65%] h-3 rounded-full' />
          <div className='flex w-full gap-3 items-center h-10'>
            <Skeleton className='size-5 rounded-full' />
            <Skeleton className='w-full h-4 rounded-full' />
          </div>
        </div>
      </div>
    </div>
  );
};

const SkeletonCardFooter2 = () => {
  return (
    <div className='flex flex-col p-4 pt-0'>
      <Divider className='mb-4 mt-0' />

      <div className='flex w-full gap-1 items-center'>
        <Skeleton className='size-5 rounded-full' />
        <Skeleton className='w-[45%] h-3 rounded-full' />
      </div>
    </div>
  );
};
