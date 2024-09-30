import { Skeleton } from '@ui/feedback/Skeleton';

export const FlowBuilderSkeleton = () => {
  return (
    <div className='flex flex-col items-center flex-1 max-h-[50%] mt-[200px]'>
      <Skeleton className='w-[300px] h-[48px] rounded-md' />
      <Skeleton className='w-[1px] h-[40px]' />
      <Skeleton className='size-2 rounded-full' />
      <Skeleton className='w-[1px] h-[40px]' />
      <Skeleton className='w-[131px] h-[48px] rounded-md' />
      <Skeleton className='w-[1px] h-[40px]' />
      <Skeleton className='size-2 rounded-full' />
      <Skeleton className='w-[1px] h-[40px]' />{' '}
      <Skeleton className='w-[300px] h-[48px] rounded-md' />
      <Skeleton className='w-[1px] h-[40px]' />
      <Skeleton className='size-2 rounded-full' />
      <Skeleton className='w-[1px] h-[40px]' />
      <Skeleton className='w-[131px] h-[48px] rounded-md' />
    </div>
  );
};
