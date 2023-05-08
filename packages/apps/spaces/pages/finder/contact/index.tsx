import type { NextPage } from 'next';
import React, { useState } from 'react';
import { PageContentLayout } from '../../../components/ui-kit/layouts';
import { SidePanel } from '../../../components/ui-kit/organisms';
import { useRouter } from 'next/router';
import { useRecoilValue } from 'recoil';
import { userData } from '../../../state';
import { FinderContact } from '../../../components/finder/finder-contact/FinderContact';
import dynamic from 'next/dynamic';
const WebChat = dynamic(() =>
  import('@openline-ai/openline-web-chat').then((res) => res.WebChat),
);
const FinderContactPage: NextPage = () => {
  const router = useRouter();
  const [isSidePanelVisible, setSidePanelVisible] = useState(false);
  const loggedInUserData = useRecoilValue(userData);

  return (
    <PageContentLayout
      isPanelOpen={isSidePanelVisible}
      isSideBarShown={router.pathname === '/'}
    >
      {router.pathname === '/finder/contact' && (
        <SidePanel
          onPanelToggle={setSidePanelVisible}
          isPanelOpen={isSidePanelVisible}
        >
          <WebChat
            apikey={`${process.env.WEB_CHAT_API_KEY}`}
            httpServerPath={`${process.env.WEB_CHAT_HTTP_PATH}`}
            wsServerPath={`${process.env.WEB_CHAT_WS_PATH}`}
            location='left'
            trackerEnabled={
              `${process.env.WEB_CHAT_TRACKER_ENABLED}` === 'true'
            }
            trackerAppId={`${process.env.WEB_CHAT_TRACKER_APP_ID}`}
            trackerId={`${process.env.WEB_CHAT_TRACKER_ID}`}
            trackerCollectorUrl={`${process.env.WEB_CHAT_TRACKER_COLLECTOR_URL}`}
            trackerBufferSize={`${process.env.WEB_CHAT_TRACKER_BUFFER_SIZE}`}
            trackerMinimumVisitLength={`${process.env.WEB_CHAT_TRACKER_MINIMUM_VISIT_LENGTH}`}
            trackerHeartbeatDelay={`${process.env.WEB_CHAT_TRACKER_HEARTBEAT_DELAY}`}
            userEmail={loggedInUserData.identity}
          />
        </SidePanel>
      )}
      <article style={{ gridArea: 'content' }}>
        <FinderContact />
      </article>
    </PageContentLayout>
  );
};

export default FinderContactPage;
