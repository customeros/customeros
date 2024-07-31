import { useNavigate } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button';
import { File04 } from '@ui/media/icons/File04';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';

import HalfCirclePattern from '../../../assets/HalfCirclePattern';

export const EmptyState = () => {
  const navigate = useNavigate();

  return (
    <div className='flex flex-col h-full w-full max-w-[448px]'>
      <div className='flex relative'>
        <FeaturedIcon
          size='lg'
          colorScheme='primary'
          className='absolute top-[26%] justify-self-center right-0 left-0'
        >
          <File04 className='size-5' />
        </FeaturedIcon>
        <HalfCirclePattern />
      </div>
      <div className='flex flex-col text-center items-center translate-y-[-200px]'>
        <p className='text-gray-700 text-md font-semibold'>
          No upcoming invoices
        </p>
        <p className='text-sm text-gray-500 my-1'>
          Schedule invoices by creating a contract with services
        </p>
        <Button
          size='sm'
          variant='outline'
          className={'mt-4 text-sm'}
          onClick={() => navigate(`?tab=account`)}
        >
          Go back
        </Button>
      </div>
    </div>
  );
};
