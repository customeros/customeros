import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

import { autorun } from 'mobx';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';

const publicPaths = ['/auth/signin', '/auth/failure'];

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
    const navigate = useNavigate();
    const location = useLocation();
    const [hidden, setHidden] = useState(false);
    const { pathname } = location;

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
          !store.isAuthenticated &&
          !store.isAuthenticating &&
          pathname === '/'
        ) {
          const hasSession = store.session.getLocalStorageSession() !== null;
          if (hasSession) {
            navigate('/organizations');
          } else {
            navigate('/auth/signin');
          }
        }

        if (store.isBootstrapping) return;

        if (!store.isAuthenticated && !publicPaths.includes(pathname)) {
          navigate('/auth/signin');
        }
        if (
          store.isAuthenticated &&
          allowedPaths.some((p) => p.startsWith(pathname))
        ) {
          navigate(
            `/organizations?preset=${store.tableViewDefs.defaultPreset}`,
          );
        }
      });

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
