import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { autorun } from 'mobx';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { validateEmail } from '@utils/email';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { Google } from '@ui/media/logos/Google';
import { useStore } from '@shared/hooks/useStore';
import { Divider } from '@ui/presentation/Divider';
import { Microsoft } from '@ui/media/logos/Microsoft';

import CustomerOsLogo from './CustomerOS-logo.png';

const providers = [
  { id: 'google', name: 'Google' },
  { id: 'azure-ad', name: 'Microsoft' },
];

export const SignIn = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [hasMagicLinkSent, setHasMagicLinkSent] = useState(false);
  const [emailValidationError, setEmailValidationError] = useState('');

  const handleSignIn = (provider: string) => {
    switch (provider) {
      case 'google':
        return store.session.authenticate('google');
      case 'azure-ad':
        return store.session.authenticate('azure-ad');

      case 'magic-link': {
        if (email.length === 0) {
          setEmailValidationError('Houston, we have a blank...');

          return;
        }

        const validationError = validateEmail(email);

        if (validationError) {
          setEmailValidationError(validationError);

          return;
        }

        setEmailValidationError('');

        return store.session.authenticate(
          'magic-link',
          { email },
          {
            onSuccess: () => {
              setHasMagicLinkSent(true);
            },
          },
        );
      }
      default:
        break;
    }
  };

  useEffect(() => {
    const dispose = autorun(() => {
      if (store.isAuthenticated) {
        navigate(`/finder?preset=${store.tableViewDefs.defaultPreset}`);
      }
    });

    return () => {
      dispose();
    };
  }, []);

  if (hasMagicLinkSent) {
    return (
      <div className='h-screen w-screen flex animate-fadeIn overflow-hidden max-h-screen max-w-screen'>
        <div className='flex-1 items-center h-screen overflow-hidden'>
          <div className='h-full flex items-center justify-center relative'>
            <div className='relative flex flex-col items-center justify-center w-[450px] h-full px-6 pb-6 bg-white'>
              <div className='h-full flex items-center justify-center relative '>
                <div className='flex flex-col items-center w-[360px]'>
                  <img
                    width={264}
                    height={264}
                    alt='CustomerOS'
                    src={CustomerOsLogo}
                  />
                  <h2 className='text-gray-900 leading-9 font-bold text-3xl py-3 mt-[-40px]'>
                    Check your email
                  </h2>
                  <p className='mb-4 text-gray-500 text-center'>
                    We've sent you an email with a magic code to {email}
                  </p>
                  <Button
                    size='md'
                    variant='outline'
                    className={cn(`w-full py-[9px] px-4`)}
                    onClick={() => {
                      setEmail('');
                      setHasMagicLinkSent(false);
                    }}
                  >
                    Try another way
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <>
      <div className='h-screen w-screen flex animate-fadeIn overflow-hidden max-h-screen max-w-screen'>
        <div className='flex-1 items-center h-screen overflow-hidden'>
          <div className='h-full flex items-center justify-center  relative'>
            <div className='flex flex-col items-center justify-center w-[450px] h-full px-6 pb-6 bg-white'>
              <img
                width={264}
                height={264}
                alt='CustomerOS'
                src={CustomerOsLogo}
              />
              <h2 className='text-gray-900 leading-9 font-bold text-3xl py-3 mt-[-40px]'>
                Welcome back
              </h2>
              <p className='text-gray-500'>Sign in to your account</p>
              {providers.map((provider, i) => {
                let icon = undefined;

                switch (provider.id) {
                  case 'google':
                    icon = <Google className='size-6' />;
                    break;
                  case 'azure-ad':
                    icon = <Microsoft className='size-6' />;
                    break;
                  default:
                    icon = undefined;
                }

                return (
                  <Button
                    size='md'
                    leftIcon={icon}
                    key={provider.id}
                    variant='outline'
                    colorScheme='gray'
                    onClick={() => handleSignIn(provider.id)}
                    isLoading={store.session.isLoading === provider.id}
                    className={cn(
                      `mt-3 w-[100%] py-[7px] px-4`,
                      i === 0 ? 'mt-6' : 'mt-3',
                    )}
                    isDisabled={
                      store.session.isLoading !== null &&
                      store.session.isLoading !== provider.id
                    }
                    rightSpinner={
                      <Spinner
                        size='sm'
                        label='Authenthicating'
                        className='text-gray-300 fill-gray-500'
                      />
                    }
                  >
                    Sign in with {provider.name}
                  </Button>
                );
              })}
              <Divider className='my-4'>
                <span className='text-sm text-gray-500 leading-none'>or</span>
              </Divider>
              <div className='flex w-full flex-col gap-4 items-center'>
                <div className='w-full'>
                  <Input
                    value={email}
                    variant='outline'
                    placeholder='Enter your email'
                    invalid={emailValidationError.length > 0}
                    onChange={(e) => setEmail(e.target.value)}
                    className='rounded-lg w-full placeholder:text-sm text-sm'
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        handleSignIn('magic-link');
                      }
                    }}
                  />
                  {emailValidationError.length > 0 && (
                    <p className='pl-[9px] text-xs text-error-500'>
                      {emailValidationError}
                    </p>
                  )}
                </div>
                <Button
                  size='md'
                  variant='outline'
                  colorScheme='primary'
                  className={cn(`w-full py-[9px] px-4`)}
                  onClick={() => handleSignIn('magic-link')}
                  isLoading={store.session.isLoading === 'magic-link'}
                  isDisabled={hasMagicLinkSent || !!store.session.isLoading}
                  rightSpinner={
                    <Spinner
                      size='sm'
                      label='Authenthicating'
                      className='text-gray-300 fill-gray-500'
                    />
                  }
                >
                  Sign in with email
                </Button>
                <p className='text-xs text-gray-500'>
                  We'll send you an email with a magic link
                </p>
              </div>

              <Divider className='my-4' />

              <div className='text-gray-500 mt-2 text-center text-xs'>
                By logging in you agree to CustomerOS&apos;s
                <div className='text-gray-500'>
                  <a
                    className='text-primary-700 mr-1 no-underline'
                    href='https://customeros.ai/legal/terms-of-service'
                  >
                    Terms of Service
                  </a>
                  <span className='mr-1'>and</span>
                  <a
                    className='text-primary-700 no-underline'
                    href='https://www.customeros.ai/legal/privacy-policy'
                  >
                    Privacy Policy
                  </a>
                  .
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
});
