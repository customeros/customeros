import React, { FC } from 'react';

import { Skeleton } from '@ui/feedback/Skeleton';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

export const TimelineItemSkeleton: FC = () => {
  return (
    <div className='mt-4 mr-6'>
      <Skeleton className='h-[0.5rem] w-[100px] rounded-md mb-4' />
      <Card className='text-sm bg-white flex flex-row max-w-[549px] shadow-xs'>
        <CardContent className='pt-5 pb-5 pl-5 pr-0 overflow-hidden flex flex-row w-full'>
          <div className='flex flex-col items-start w-full justify-between h-full'>
            <Skeleton className='w-[33%] h-[0.75rem] rounded-md' />
            <Skeleton className='w-[95%] h-[0.5rem] rounded-md' />
            <Skeleton className='w-[95%] h-[0.5rem] rounded-md' />
            <Skeleton className='w-[95%] h-[0.5rem] rounded-md' />
          </div>
        </CardContent>
        <CardFooter className='pt-5 pb-5 pr-5 pl-0 ml-1'>
          <Skeleton className='h-[70px] w-[54px] rounded-md' />
        </CardFooter>
      </Card>
    </div>
  );
};
