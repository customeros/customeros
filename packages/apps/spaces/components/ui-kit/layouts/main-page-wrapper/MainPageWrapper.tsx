import '@openline-ai/openline-web-chat/dist/esm/index.css';
import { Configuration, FrontendApi, Session } from '@ory/client';
import { edgeConfig } from '@ory/integrations/next';
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { getUserName } from '../../../../utils';
import client from '../../../../apollo-client';
import { ApolloProvider } from '@apollo/client';
import { logoutUrlState, userData } from '../../../../state';
import { useSetRecoilState } from 'recoil';
import { WebRTCContextProvider } from '../../../../context';
import { WebRTCCallProgress, WebRTCInboundNotification } from '../../molecules';

const ory = new FrontendApi(new Configuration(edgeConfig));

export const MainPageWrapper = ({ children }: any) => {
  const router = useRouter();
  // const setTheme = (theme) => {
  //     document.documentElement.className = theme;
  //     localStorage.setItem('theme', theme);
  // }
  // const getTheme = () => {
  //     const theme = localStorage.getItem('theme');
  //     theme && setTheme(theme);
  // }
  const setLogoutUrl = useSetRecoilState(logoutUrlState);

  const [session, setSession] = useState<Session | undefined>();
  const setUserEmail = useSetRecoilState(userData);

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
        // @ts-expect-error analytics is added to window object from script
        if (window?.analytics && process.env.NODE_ENV === 'production') {
          // @ts-expect-error analytics is added to window object from script
          window.analytics.identify(data.id, {
            email: getUserName(data.identity),
          });
          // @ts-expect-error analytics is added to window object from script
          window.analytics.track('Signed in', {
            // @ts-expect-error analytics is added to window object from script
            browser: data.devices?.user_agent || 'unknown',
          });
        }

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
