import { Button } from '@ui/form/Button/Button';
import HalfCirclePattern from '@shared/assets/HalfCirclePattern';

import Illustration from './graphic.svg?react';

interface EmptyStateProps {
  onClick: () => void;
}

export const EmptyState = ({ onClick }: EmptyStateProps) => {
  return (
    <div className='flex justify-center h-full bg-white'>
      <div className='flex flex-col h-[500px] w-[500px]'>
        <div className='flex relative'>
          <Illustration className='w-[152px] h-[120px] absolute top-[25%] right-[35%]' />
          <HalfCirclePattern width={500} height={500} />
        </div>
        <div className='flex flex-col text-center items-center top-[5vh] transform translate-y-[-230px]'>
          <p className='text-gray-900 text-md font-semibold'>
            Agents orchestrate your business
          </p>
          <p className='max-w-[400px] text-sm text-gray-600 my-1'>
            Agents let you manage your whole business from <br /> identifying
            website visitors, to qualifying them to <br /> sending outbound
            messages.
          </p>

          <Button
            variant='outline'
            onClick={onClick}
            colorScheme='primary'
            className='mt-4 min-w-min text-sm bg-white'
          >
            Get started
          </Button>
        </div>
      </div>
    </div>
  );
};
