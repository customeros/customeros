import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { autorun } from 'mobx';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { Google } from '@ui/media/logos/Google';
import { useStore } from '@shared/hooks/useStore';
import { Divider } from '@ui/presentation/Divider';
import { Microsoft } from '@ui/media/logos/Microsoft';

import Background from './login-bg.png';
import CustomerOsLogo from './CustomerOS-logo.png';

import BackgroundGridDot from '/backgrounds/grid/backgroundGridDot.png';

const providers = [
  { id: 'google', name: 'Google' },
  { id: 'azure-ad', name: 'Microsoft' },
];

export const SignIn = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [hasMagicLinkSent, setHasMagicLinkSent] = useState(false);

  const handleSignIn = (provider: string) => {
    switch (provider) {
      case 'google':
        return store.session.authenticate('google');
      case 'azure-ad':
        return store.session.authenticate('azure-ad');

      case 'magic-link':
        return store.session.authenticate(
          'magic-link',
          { email },
          {
            onSuccess: () => {
              setEmail('');
              setHasMagicLinkSent(true);
            },
          },
        );
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
      <div className='h-screen w-screen flex animate-fadeIn'>
        <div className='flex-1'>
          <div className='h-[50%] w-[100%]'>
            <img
              alt=''
              src={BackgroundGridDot}
              className='top-[-10%] relative w-[480px] m-auto'
            />
          </div>
          <div className='h-full flex items-center justify-center relative top-[-50%]'>
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
              <p className='text-gray-500'>
                We've sent you an email with a magic code
              </p>
              <Divider className='my-4' />
              <Button
                size='md'
                variant='outline'
                colorScheme='primary'
                className={cn(`w-full py-[9px] px-4`)}
                onClick={() => setHasMagicLinkSent(false)}
              >
                Resend email
              </Button>
            </div>
          </div>
        </div>
        <img
          src={Background}
          alt='Background'
          className=' flex-1 bg-cover rounded-s-[80px] bg-no-repeat h-full w-[50vw]'
        />
      </div>
    );
  }

  return (
    <>
      <div className='h-screen w-screen flex animate-fadeIn'>
        <div className='flex-1'>
          <div className='h-[50%] w-[100%]'>
            <img
              alt=''
              src={BackgroundGridDot}
              className='top-[-10%] relative w-[480px] m-auto'
            />
          </div>
          <div className='h-full flex items-center justify-center relative top-[-55%]'>
            <div className='flex flex-col items-center w-[360px]'>
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
                <Input
                  value={email}
                  variant='outline'
                  placeholder='Enter your email'
                  onChange={(e) => setEmail(e.target.value)}
                  className='rounded-lg w-full placeholder:text-sm text-sm'
                />
                <Button
                  size='md'
                  variant='outline'
                  colorScheme='primary'
                  className={cn(`w-full py-[9px] px-4`)}
                  isLoading={store.session.isLoading === 'magic-link'}
                  isDisabled={hasMagicLinkSent || !!store.session.isLoading}
                  onClick={() => email.length > 0 && handleSignIn('magic-link')}
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
        <img
          src={Background}
          alt='Background'
          className=' flex-1 bg-cover rounded-s-[80px] bg-no-repeat h-full w-[50vw]'
        />
      </div>
    </>
  );
});
