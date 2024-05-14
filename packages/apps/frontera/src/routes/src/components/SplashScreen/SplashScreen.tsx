import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';

const publicPahts = ['/auth/signin', '/auth/failure', '/auth/success'];

const allowedPaths = [
  '/auth',
  '/organizations/',
  '/organization',
  '/invoices',
  '/renewals',
  '/customer-map',
  '/settings',
  '/prospects',
];

export const SplashScreen = observer(
  ({ children }: { children: React.ReactNode }) => {
    const store = useStore();
    const [hidden, setHidden] = useState(false);

    const showSplash =
      !store.isBootstrapped && !publicPahts.includes(window.location.pathname);
    const hide = hidden || publicPahts.includes(window.location.pathname);
    const render =
      publicPahts.some((p) => p.startsWith(window.location.pathname)) ||
      store.isBootstrapped;

    useEffect(() => {
      if (
        store.isBootstrapped ||
        publicPahts.includes(window.location.pathname)
      ) {
        setTimeout(() => {
          setHidden(true);
        }, 500);
      }
    }, [store.isBootstrapped]);

    useEffect(() => {
      if (!allowedPaths.some((path) => location.pathname.startsWith(path))) {
        window.location.pathname = '/auth/signin';
      }
    }, []);

    return (
      <>
        {render && children}
        <div
          className={cn(
            'absolute flex items-center justify-center top-0 right-0 bottom-0 left-0 z-10 bg-white opacity-0 duration-500 transition-opacity',
            showSplash && 'opacity-100',
            hide && 'hidden',
          )}
        >
          <Spinner label='loading' className='text-gray-300 fill-gray-500' />
        </div>
      </>
    );
  },
);
