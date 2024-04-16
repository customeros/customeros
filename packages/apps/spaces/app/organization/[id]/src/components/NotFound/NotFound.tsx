'use client';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon2';

import HalfCirclePattern from '../../../../../src/assets/HalfCirclePattern';

export default function NotFound() {
  return (
    <div className='p-0 flex-1 flex-col bg-no-repeat bg-contain w-[100vw] relative items-center justify-center bg-gray-25 border-1 border-gray-200 rounded-xl'>
      <div
        className='absolute h-[50vh] max-h-[768px] w-[768px] top-[50%] left-[50%]'
        style={{ transform: 'translate(-50%, -90%) rotate(180deg)' }}
      >
        <HalfCirclePattern />
      </div>
      <div className='flex relative flex-col items-center justify-center h-[50vh]'>
        <FeaturedIcon colorScheme='primary' size='lg'>
          <SearchSm boxSize='5' />
        </FeaturedIcon>
        <h2 className='text-5xl text-gray-900 py-6 font-semibold'>
          This organization cannot be found
        </h2>
        <p className='text-gray-600 text-2xl pb-12 px-8 text-center'>
          It appears the organization does not exist or you do not have
          sufficient rights to preview it.
        </p>
      </div>
    </div>
  );
}
