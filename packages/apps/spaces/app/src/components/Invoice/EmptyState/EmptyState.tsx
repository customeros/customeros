import { useRouter } from 'next/navigation';

import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { File04 } from '@ui/media/icons/File04';

import HalfCirclePattern from '../../../assets/HalfCirclePattern';

export const EmptyState = () => {
  const router = useRouter();

  return (
    <div className='flex flex-col h-full w-full max-w-[448px]'>
      <div className='flex relative'>
        <FeaturedIcon
          size='lg'
          colorScheme='primary'
          className='absolute top-[23%] right-[45%]'
        >
          <File04 boxSize='5' />
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
          variant='outline'
          size='sm'
          className={'mt-4 text-sm'}
          onClick={() => router.push(`?tab=account`)}
        >
          Go back
        </Button>
      </div>
    </div>
  );
};
