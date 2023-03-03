import '@openline-ai/openline-web-chat/dist/esm/index.css';
import { Configuration, FrontendApi, Session } from '@ory/client';
import { edgeConfig } from '@ory/integrations/next';
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { getUserName } from '../../../../utils';
import client from '../../../../apollo-client';
import { ApolloProvider } from '@apollo/client';
import { useRecoilState, useSetRecoilState } from 'recoil';
import { logoutUrlState, userData } from '../../../../state';
import { useSetRecoilState } from 'recoil';
import { WebRTCContextProvider } from '../../../../context';

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

  useEffect(() => {
    if (router.asPath.startsWith('/login')) {
      return;
    }
    ory
      .toSession()
      .then(({ data }) => {
        // User has a session!
        setSession(data);
        setUserEmail({ identity: getUserName(data.identity) });
        // Create a logout url
        ory.createBrowserLogoutFlow().then(({ data }) => {
          setLogoutUrl(data.logout_url);
        });
      })
      .catch(() => {
        // Redirect to login page
        return router.push(edgeConfig.basePath + '/ui/login');
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
      <WebRTCContextProvider username={userEmail}> {children}</WebRTCContextProvider>
    </ApolloProvider>
  );
};
