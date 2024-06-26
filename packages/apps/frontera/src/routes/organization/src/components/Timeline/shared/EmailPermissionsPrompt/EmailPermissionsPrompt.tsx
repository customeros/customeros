import { FC } from 'react';
import { useNavigate } from 'react-router-dom';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Google } from '@ui/media/logos/Google';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';

export const MissingPermissionsPrompt: FC<{
  modal: boolean;
}> = ({ modal }) => {
  const navigate = useNavigate();

  return (
    <form
      className={cn(
        modal
          ? 'bg-grayBlue-50 border-t w-full border-dashed border-gray-200 max-h-[50vh]'
          : 'bg-white rounded-lg max-h-[auto] ',
        'flex items-center mt-4 p-6 overflow-visible rounded-b-2xl justify-center',
      )}
    >
      <div
        className={cn(
          modal ? 'bg-[#F8F9FC]' : 'bg-white',
          'flex flex-col items-center p-6',
        )}
      >
        <FeaturedIcon size='lg' className='mb-4' colorScheme='gray'>
          <Mail01 className='text-gray-700 size-6' />
        </FeaturedIcon>
        <p className='text-gray-700 font-semibold mb-1'>
          Allow CustomerOS to send emails
        </p>

        <p className='text-gray-500 mb-6 text-center'>
          To send emails, you need to allow CustomerOS to connect to your email
          account
        </p>
        <Button
          variant='outline'
          colorScheme='gray'
          onClick={() => {
            navigate('/settings');
          }}
        >
          <Google className='mr-2' />
          Allow with Google or Microsoft
        </Button>
      </div>
    </form>
  );
};
