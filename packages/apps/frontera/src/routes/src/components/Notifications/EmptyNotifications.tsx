import React from 'react';

import { Lotus } from '@ui/media/icons/Lotus';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';

import HalfCirclePattern from '../../assets/HalfCirclePattern';

export const EmptyNotifications = () => {
  return (
    <div className='relative flex flex-col items-center max-w-[448px] px-4 py-1 mt-5 overflow-hidden'>
      <div className='absolute h-[400px] w-[448px] transform -translate-y-[75px]'>
        <HalfCirclePattern />
      </div>
      <FeaturedIcon className='mt-5'>
        <Lotus />
      </FeaturedIcon>
      <h1 className='mt-8 mb-1 text-4 z-10 font-semibold leading-6 text-gray-700'>
        No notifications for now
      </h1>
      <span className='text-center z-10 text-sm leading-5 text-gray-500'>
        Enjoy the quiet moment. Explore other corners of the app or take a deep
        breath and savor the calm.
      </span>
    </div>
  );
};
