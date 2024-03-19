'use client';
import React from 'react';
import Image from 'next/image';

import { BuiltInProviderType } from 'next-auth/providers';
import { signIn, LiteralUnion, ClientSafeProvider } from 'next-auth/react';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { Google } from '@ui/media/logos/Google';
import { Microsoft } from '@ui/media/logos/Microsoft';

import Background from './login-bg.png';
import CustomOsLogo from './CustomerOS-logo.png';
import BackgroundGridDot from '../../../public/backgrounds/grid/backgroundGridDot.png';

export default function SignIn({
  providers,
}: {
  providers: Record<
    LiteralUnion<BuiltInProviderType, string>,
    ClientSafeProvider
  > | null;
}) {
  return (
    <>
      <div className='h-screen'>
        <div className='h-[50%] w-[100%]'>
          <Image
            alt=''
            src={BackgroundGridDot}
            className='top-[-10%] relative w-[480px] m-auto'
          />
        </div>
        <div className='h-full flex items-center justify-center relative top-[-50%]'>
          <div className='flex flex-col items-center w-[360px]'>
            <Image
              src={CustomOsLogo}
              alt='CustomerOS'
              className='size-[264px]'
            />
            <h2 className='text-gray-900 leading-9 font-bold text-3xl py-3 mt-[-40px]'>
              Welcome back
            </h2>
            <p className='text-gray-500'>Sign in to your account</p>
            {providers &&
              Object.values(providers).map((provider, i) => {
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
                    key={provider.name}
                    size='md'
                    variant='outline'
                    colorScheme='gray'
                    leftIcon={icon}
                    onClick={() => signIn(provider.id)}
                    className={cn(
                      ` mt-3 w-[100%] py-[7px] px-4`,
                      i === 0 ? 'mt-6' : 'mt-3',
                    )}
                  >
                    Sign in with {provider.name}
                  </Button>
                );
              })}
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
      <div
        className=' bg-cover rounded-s-[80px] bg-no-repeat'
        style={{ backgroundImage: `url(${Background.src})` }}
      />
    </>
  );
}
