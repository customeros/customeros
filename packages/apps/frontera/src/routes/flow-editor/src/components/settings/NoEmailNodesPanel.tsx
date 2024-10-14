import React from 'react';

import { FeaturedIcon } from '@ui/media/Icon';
import { Sliders04 } from '@ui/media/icons/Sliders04';

export const NoEmailNodesPanel = () => {
  return (
    <div
      className='flex flex-col h-[400px] bg-no-repeat bg-cover '
      style={{
        backgroundImage: `url(/backgrounds/organization/dotted-bg-pattern.svg)`,
        backgroundPositionX: 'center',
      }}
    >
      <h2 className='text-base font-medium'>Flow Settings</h2>

      <div className='flex flex-col items-center mt-8'>
        <FeaturedIcon size='lg' colorScheme='gray'>
          <Sliders04 className='text-primary-600' />
        </FeaturedIcon>
        <span className='text-gray-700 font-semibold mt-8'>
          Step up your flow game
        </span>
        <span className='text-gray-500 mt-1 mb-6 text-sm text-center'>
          To configure your flow’s settings, start by adding some steps to it
          first
        </span>
      </div>
    </div>
  );
};
