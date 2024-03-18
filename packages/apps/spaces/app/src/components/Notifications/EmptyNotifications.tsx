import React from 'react';

import { Lotus } from '@ui/media/icons/Lotus';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon2';

import HalfCirclePattern from '../../assets/HalfCirclePattern';

export const EmptyNotifications = () => {
  return (
    <div className='relative flex flex-col items-center max-w-[448px] px-4 py-1 mt-5 overflow-hidden'>
      <div className='absolute h-[400px] w-[448px]'>
        <HalfCirclePattern />
      </div>
      <FeaturedIcon className='mt-[20px]' colorScheme='primary'>
        <Lotus />
      </FeaturedIcon>
      <div className='mt-4 mb-1 text-[16px] font-bold leading-5 text-gray-900'>
        No notifications for now
      </div>
      <div className='text-center text-[14px] leading-5 text-gray-500'>
        Enjoy the quiet moment. Explore other corners of the app or take a deep
        breath and savor the calm.
      </div>
    </div>
  );
};
