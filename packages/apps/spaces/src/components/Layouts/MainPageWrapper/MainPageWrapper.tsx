import '@openline-ai/openline-web-chat/dist/esm/index.css';
import { Configuration, FrontendApi, Session } from '@ory/client';
import { edgeConfig } from '@ory/integrations/next';
import React, { useEffect, useState } from 'react';
import { WebChat } from '@openline-ai/openline-web-chat';
import { useRouter } from 'next/router';
import { getUserName } from '../../../utils';

const ory = new FrontendApi(new Configuration(edgeConfig));

export const MainPageWrapper = ({ children }: any) => {
  const router = useRouter();
  const [isSidePanelVisible, setSidePanelVisible] = useState(false);
  // const setTheme = (theme) => {
  //     document.documentElement.className = theme;
  //     localStorage.setItem('theme', theme);
  // }
  // const getTheme = () => {
  //     const theme = localStorage.getItem('theme');
  //     theme && setTheme(theme);
  // }

  const [session, setSession] = useState<Session | undefined>();
  const [userEmail, setUserEmail] = useState<string | undefined>();
  const [logoutUrl, setLogoutUrl] = useState<string | undefined>();

  useEffect(() => {
    if (router.asPath.startsWith('/login')) {
      return;
    }
    ory
      .toSession()
      .then(({ data }) => {
        // User has a session!
        setSession(data);
        setUserEmail(getUserName(data.identity));
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
    <>
      {/*{router.pathname === '/' && (*/}
      {/*    <SidePanel*/}
      {/*        userEmail={userEmail}*/}
      {/*        logoutUrl={logoutUrl}*/}
      {/*        isOpen={isSidePanelVisible}*/}
      {/*        onOpen={() => setSidePanelVisible(true)}*/}
      {/*        onClose={() => setSidePanelVisible(false)}*/}
      {/*    />*/}
      {/*)}*/}

      {children}

      <WebChat
        apikey={`${process.env.WEB_CHAT_API_KEY}`}
        httpServerPath={`${process.env.WEB_CHAT_HTTP_PATH}`}
        wsServerPath={`${process.env.WEB_CHAT_WS_PATH}`}
        location='left'
        trackerEnabled={`${process.env.WEB_CHAT_TRACKER_ENABLED}` === 'true'}
        trackerAppId={`${process.env.WEB_CHAT_TRACKER_APP_ID}`}
        trackerId={`${process.env.WEB_CHAT_TRACKER_ID}`}
        trackerCollectorUrl={`${process.env.WEB_CHAT_TRACKER_COLLECTOR_URL}`}
        trackerBufferSize={`${process.env.WEB_CHAT_TRACKER_BUFFER_SIZE}`}
        trackerMinimumVisitLength={`${process.env.WEB_CHAT_TRACKER_MINIMUM_VISIT_LENGTH}`}
        trackerHeartbeatDelay={`${process.env.WEB_CHAT_TRACKER_HEARTBEAT_DELAY}`}
        userEmail={userEmail || ''}
      />
    </>
  );
};
