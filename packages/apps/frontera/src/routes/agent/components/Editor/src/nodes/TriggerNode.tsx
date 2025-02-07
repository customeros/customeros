import { Icon } from '@ui/media/Icon';

import { Handle } from '../components';

export const TriggerNode = () => {
  return (
    <>
      <div
        className={`h-[83px] w-[300px] bg-white border border-grayModern-300 rounded-lg group relative cursor-pointer flex flex-col items-center`}
      >
        <div
          data-test={'flow-trigger-block'}
          className='px-4 bg-gray-25 text-xs h-full flex items-center w-full rounded-t-lg justify-center border-b border-dashed border-gray-300 text-gray-500'
        >
          Flow triggers when
        </div>

        <div className='flex items-center justify-between w-full p-4  h-[56px]'>
          <div className='truncate text-sm flex items-center'>
            <div className='size-6 mr-2 bg-gray-50 border border-gray-100 rounded flex items-center justify-center'>
              <Icon name='user-plus-01' className='text-gray-500' />
            </div>

            <span className='font-medium '>People are added to this flow</span>
          </div>
        </div>
        <Handle type='target' />
        <Handle type='source' />
      </div>
    </>
  );
};
