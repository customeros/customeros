import React, { useEffect } from 'react';
import { useRouter } from 'next/router';
import client from '../../../../apollo-client';
import { ApolloProvider } from '@apollo/client';
import { useJune } from '@spaces/hooks/useJune';
import { useSession } from 'next-auth/react';

export const MainPageWrapper = ({ children }: any) => {
  const router = useRouter();
  const analytics = useJune();
  const { data: session } = useSession();


  useEffect(() => {
    if (analytics && session) {
      analytics.user().then((user) => {
        if (!user || user.id() === null) {
          analytics?.identify(session.user?.email, {
            email: session.user?.email,
          });
        }
      });
      analytics.pageView(router.asPath);
    }
  }, [session, analytics]);

  return (
    <ApolloProvider client={client}>
      {/*<WebRTCContextProvider>*/}
      {/*  <WebRTCInboundNotification />*/}
      {/*  <WebRTCCallProgress />*/}
      {children}
      {/*</WebRTCContextProvider>*/}
    </ApolloProvider>
  );
};
