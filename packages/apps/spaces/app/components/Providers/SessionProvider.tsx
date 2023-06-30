'use client';

import { useEffect, useState } from 'react';
import { useSetRecoilState } from 'recoil';
import { edgeConfig } from '@ory/integrations/next';
import { useRouter, usePathname } from 'next/navigation';
import { Configuration, FrontendApi, Session } from '@ory/client';

import { useJune } from '@spaces/hooks/useJune';
import { userData } from '@spaces/globalState/userData';
import { getUserName } from '@spaces/utils/getLoggedInData';

const ory = new FrontendApi(new Configuration(edgeConfig));

const getReturnToUrl = (): string => {
  if (window.location.origin.startsWith('http://localhost')) {
    return '';
  }
  return '?return_to=' + window.location.origin;
};

export const SessionProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const router = useRouter();
  const analytics = useJune();
  const pathname = usePathname();

  const [session, setSession] = useState<Session | undefined>();
  const setUserData = useSetRecoilState(userData);

  useEffect(() => {
    if (pathname?.startsWith('/login')) {
      return;
    }

    (async () => {
      try {
        const { data: sessionData } = await ory.toSession();
        setSession(sessionData);
        sessionStorage.setItem('session', JSON.stringify(sessionData));
        setUserData({
          identity: getUserName(sessionData.identity),
          id: sessionData.id,
        });

        const {
          data: { logout_url },
        } = await ory.createBrowserLogoutFlow();
        sessionStorage.setItem('logout_url', logout_url);
      } catch (e) {
        return router.push(
          edgeConfig.basePath + '/ui/login' + getReturnToUrl(),
        );
      }
    })();
  }, [pathname, router, setUserData]);

  useEffect(() => {
    if (analytics && session) {
      analytics.user().then((user) => {
        if (!user || user.id() === null) {
          analytics?.identify(session.identity.traits.email, {
            email: session.identity.traits.email,
          });
        }
      });
      if (pathname) {
        analytics.pageView(pathname);
      }
    }
  }, [session, analytics, pathname]);

  if (pathname && !session) {
    if (pathname.startsWith('/login')) {
      return <>{children}</>;
    }
    if (pathname !== '/login') {
      // TODO: this should return a skeleton
      return null;
    }
  }

  return <>{children}</>;
};
