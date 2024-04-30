import { Skeleton } from '@ui/feedback/Skeleton/Skeleton';
import { Card, CardContent } from '@ui/presentation/Card/Card';

export const MessageCardSkeleton = () => {
  return (
    <Card className='text-sm bg-white flex shadow-xs border border-gray-200 w-full'>
      <CardContent className='p-3 w-full'>
        <div className='flex gap-4'>
          <Skeleton className='w-10 h-10 rounded-[6px]' />
          <div className='flex flex-col gap-2 w-full h-full'>
            <Skeleton className='h-[16px] w-[25%]' />
            <Skeleton className='h-[12px] w-[50%]' />
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
