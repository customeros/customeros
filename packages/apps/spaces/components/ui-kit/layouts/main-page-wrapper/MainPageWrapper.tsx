import { Configuration, FrontendApi, Session } from '@ory/client';
import { edgeConfig } from '@ory/integrations/next';
import dynamic from 'next/dynamic';
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { getUserName } from '../../../../utils';
import client from '../../../../apollo-client';
import { ApolloProvider } from '@apollo/client';
import { logoutUrlState, userData } from '../../../../state';
import { useSetRecoilState } from 'recoil';
import { WebRTCContextProvider } from '../../../../context';
import { useJune } from '@spaces/hooks/useJune';

const WebRTCInboundNotification = dynamic(
  () =>
    import('../../molecules/webRTC/WebRTCInboundNotification').then(
      (res) => res.WebRTCInboundNotification,
    ),
  { ssr: true },
);

const WebRTCCallProgress = dynamic(
  () =>
    import('../../molecules/webRTC/WebRTCCallProgress').then(
      (res) => res.WebRTCCallProgress,
    ),
  { ssr: true },
);

const ory = new FrontendApi(new Configuration(edgeConfig));

export const MainPageWrapper = ({ children }: any) => {
  const router = useRouter();
  const analytics = useJune();

  const setLogoutUrl = useSetRecoilState(logoutUrlState);

  const [session, setSession] = useState<Session | undefined>();
  const setUserEmail = useSetRecoilState(userData);

  useEffect(() => {
    if (analytics && session) {
      analytics.user().then((user) => {
        if (!user || user.id() === null) {
          analytics?.identify(session.identity.traits.email, {
            email: session.identity.traits.email,
          });
        }
      });
      analytics.pageView(router.asPath);
    }
  }, [session, analytics]);

  const getReturnToUrl: () => string = () => {
    if (window.location.origin.startsWith('http://localhost')) {
      return '';
    }
    return '?return_to=' + window.location.origin;
  };

  useEffect(() => {
    if (router.asPath.startsWith('/login')) {
      return;
    }
    ory
      .toSession()
      .then(({ data }) => {
        // User has a session!
        setSession(data);
        setUserEmail({ identity: getUserName(data.identity), id: data.id });

        // Create a logout url
        ory.createBrowserLogoutFlow().then(({ data }) => {
          setLogoutUrl(data.logout_url);
        });
      })
      .catch(() => {
        // Redirect to login page
        return router.push(
          edgeConfig.basePath + '/ui/login' + getReturnToUrl(),
        );
      });
  }, [router]);

  if (!session) {
    if (router.asPath.startsWith('/login')) {
      return <>{children}</>;
    }
    if (router.asPath !== '/login') {
      return null;
    }
  }

  return (
    <ApolloProvider client={client}>
      <WebRTCContextProvider>
        <WebRTCInboundNotification />
        <WebRTCCallProgress />
        {children}
      </WebRTCContextProvider>
    </ApolloProvider>
  );
};
