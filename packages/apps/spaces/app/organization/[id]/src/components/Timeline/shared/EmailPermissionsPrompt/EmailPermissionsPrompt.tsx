import { FC } from 'react';

import { signIn } from 'next-auth/react';

import { cn } from '@ui/utils/cn';
import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Google } from '@ui/media/logos/Google';
import { toastError } from '@ui/presentation/Toast';

export const MissingPermissionsPrompt: FC<{
  modal: boolean;
}> = ({ modal }) => {
  const signInWithScopes = async () => {
    const scopes = [
      'openid',
      'email',
      'profile',
      'https://www.googleapis.com/auth/gmail.readonly',
      'https://www.googleapis.com/auth/gmail.send',
      'https://www.googleapis.com/auth/calendar.readonly',
    ];

    try {
      await signIn(
        'google',
        { callbackUrl: window.location.href },
        {
          prompt: 'login',
          scope: scopes.join(' '),
        },
      );
    } catch (error) {
      toastError('Something went wrong!', `unable-to-sign-in-with-scopes`);
    }
  };

  return (
    <form
      className={cn(
        modal
          ? 'bg-[#F8F9FC] border-t border-dashed border-gray-200 max-h-[50vh]'
          : 'bg-white rounded-lg max-h-[auto] ',
        'flex items-center mt-4 p-6 overflow-visible rounded-b-full',
      )}
    >
      <div
        className={cn(
          modal ? 'bg-[#F8F9FC]' : 'bg-white',
          'flex flex-col items-center p-6',
        )}
      >
        <FeaturedIcon size='md' minW='10' colorScheme='gray' mb={4}>
          <Mail01 color='gray.700' boxSize='6' />
        </FeaturedIcon>
        <p className='text-gray-700 font-semibold mb-1'>
          Allow CustomerOS to send emails
        </p>

        <p className='text-gray-500 mb-6 text-center'>
          To send emails, you need to allow CustomerOS to connect to your gmail
          account
        </p>
        <Button variant='outline' colorScheme='gray' onClick={signInWithScopes}>
          <Google mr={2} />
          Allow with google
        </Button>
      </div>
    </form>
  );
};
