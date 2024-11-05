import { cn } from '@ui/utils/cn';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';

export const SendConnectionRequestActionNode = () => {
  return (
    <>
      <div className='text-sm flex items-center justify-between overflow-hidden w-full'>
        <div className='truncate text-sm flex items-center'>
          <div
            className={cn(
              `size-6 min-w-6 mr-2 bg-blue-50 text-blue-500 border border-blue-100 rounded flex items-center justify-center`,
            )}
          >
            <LinkedinOutline className='size-4 text-blue-500' />
          </div>

          <span className='truncate'>Send connection request</span>
        </div>
      </div>
    </>
  );
};
