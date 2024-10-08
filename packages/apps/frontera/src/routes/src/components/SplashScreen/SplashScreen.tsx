import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

import { autorun } from 'mobx';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { LoadingScreen } from '@shared/components/SplashScreen/components';

// `/auth/success` is omitted from the list of public paths so that the spinner continues to show after a successful login
// while the user is redirected to the organizations page and bootstrapping is still in progress
const publicPaths = ['/auth/signin', '/auth/failure'];
const privatePaths = [
  '/finder',
  '/organization/',
  '/settings',
  '/invoices',
  '/flow-editor',
  '/renewals',
  '/customer-map',
];

export const SplashScreen = observer(
  ({ children }: { children: React.ReactNode }) => {
    const store = useStore();
    const navigate = useNavigate();
    const location = useLocation();
    const [hidden, setHidden] = useState(false);
    const { pathname } = location;

    const showSplash = !store.isBootstrapped && !publicPaths.includes(pathname);
    const render =
      publicPaths.some((p) => p.startsWith(pathname)) || store.isBootstrapped;

    useEffect(() => {
      if (store.isBootstrapped || publicPaths.includes(pathname)) {
        setTimeout(() => {
          setHidden(true);
        }, 500);
      }
    }, [store.isBootstrapped, store.demoMode, pathname]);

    useEffect(() => {
      if (store.demoMode) return;

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
              navigate('/finder');
            }
          }
        }
      });

      return () => dispose();
    }, []);

    if (
      store.demoMode ||
      [...publicPaths, ...privatePaths, '/auth/success'].every(
        (v) => !pathname.startsWith(v),
      )
    ) {
      return children;
    }

    return (
      <>
        {render && children}

        <LoadingScreen
          hide={hidden}
          showSplash={showSplash}
          isLoaded={store.isBootstrapped || publicPaths.includes(pathname)}
        />
      </>
    );
  },
);
