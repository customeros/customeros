import { useNavigate } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import HalfCirclePattern from '@shared/assets/HalfCirclePattern';

export const Error = () => {
  const navigate = useNavigate();

  return (
    <div className='flex-1 flex flex-col bg-no-repeat bg-contain h-screen w-screen relative items-center justify-center'>
      <div className='absolute h-[50vh] max-h-[768px] w-[768px] top-[50%] left-[50%] transform -translate-x-[50%] -translate-y-[90%] rotate-180'>
        <HalfCirclePattern />
      </div>
      <div className='relative flex flex-col items-center justify-center h-1/2'>
        <FeaturedIcon colorScheme='primary' size='lg'>
          <SearchSm className='size-5' />
        </FeaturedIcon>
        <h2 className='font-semibold text-6xl leading-[80px] text-gray-900 py-6'>
          Something went wrong
        </h2>
        <p className='text-gray-600 text-xl pb-12 px-8 leading-[30px]'>
          There was a small hiccup in the success plan. Let’s get you back to a
          familiar place.
        </p>
        <Button
          colorScheme='primary'
          variant='outline'
          size='lg'
          onClick={() => navigate(-1)}
        >
          Try again
        </Button>
      </div>
    </div>
  );
};
