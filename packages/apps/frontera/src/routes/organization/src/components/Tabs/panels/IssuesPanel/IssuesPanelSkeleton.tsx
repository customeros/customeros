import { Skeleton } from '@ui/feedback/Skeleton';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

export const IssuesPanelSkeleton = () => {
  return (
    <OrganizationPanel title='Issues'>
      <div className='flex w-full justify-start'>
        <Skeleton className='rounded-lg w-[50px] mb-1 h-[16px]' />
      </div>
      {Array.from({ length: 3 }).map((_, i) => (
        <Card
          className='w-full shadow-xs cursor-pointer bg-white border-1 border-gray-200 rounded-lg p-3'
          key={i}
        >
          <CardContent className='p-0'>
            <div className='flex flex-1 gap-4 items-start flex-wrap'>
              <Skeleton className='rounded-full size-10' />

              <div className='flex flex-col flex-1'>
                <div className='flex justify-between'>
                  <Skeleton className='rounded-full h-3 w-[50%] mb-2' />
                </div>

                <Skeleton className='rounded-full h-3 w-[55%] mb-2' />
                <Skeleton className='rounded-full h-3 w-[45%]' />
              </div>
              <Skeleton className='block static rounded-md h-6 w-10' />
            </div>
          </CardContent>
        </Card>
      ))}
    </OrganizationPanel>
  );
};
