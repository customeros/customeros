import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

import { when, autorun } from 'mobx';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';

// `/auth/success` is omitted from the list of public paths so that the spinner continues to show after a successful login
// while the user is redirected to the organizations page and bootstrapping is still in progress
const publicPaths = ['/auth/signin', '/auth/failure'];
const privatePaths = [
  '/organizations',
  '/organization/',
  '/settings',
  '/invoices',
  '/renewals',
  '/customer-map',
];

export const SplashScreen = observer(
  ({ children }: { children: React.ReactNode }) => {
    const store = useStore();
    const navigate = useNavigate();
    const location = useLocation();
    const [hidden, setHidden] = useState(false);
    const { pathname, search } = location;

    const showSplash = !store.isBootstrapped && !publicPaths.includes(pathname);
    const hide = hidden || publicPaths.includes(pathname);
    const render =
      publicPaths.some((p) => p.startsWith(pathname)) || store.isBootstrapped;

    useEffect(() => {
      if (store.isBootstrapped || publicPaths.includes(pathname)) {
        setTimeout(() => {
          setHidden(true);
        }, 500);
      }
    }, [store.isBootstrapped, pathname]);

    useEffect(() => {
      const dispose = autorun(() => {
        if (
          store.session.isBootstrapped &&
          !store.session.isBootstrapping &&
          (pathname === '/' || privatePaths.some((p) => pathname.startsWith(p)))
        ) {
          if (!store.session.isAuthenticated) {
            navigate('/auth/signin');
          } else {
            if (pathname === '/') {
              navigate('/organizations');
            }
          }
        }
      });

      when(
        () =>
          store.isAuthenticated &&
          typeof store.tableViewDefs.defaultPreset !== 'undefined' &&
          pathname.startsWith('/organizations') &&
          !search.includes('preset'),
        () => {
          navigate(
            `/organizations?preset=${store.tableViewDefs.defaultPreset}`,
          );
        },
      );

      return () => dispose();
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
