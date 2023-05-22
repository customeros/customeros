import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode, useState } from 'react';
import classNames from 'classnames';
import styles from './page-content-layout.module.scss';
import { SidePanel } from '@spaces/organisms/side-panel';

interface PageContentLayout {
  isSideBarShown: boolean;
  children: ReactNode;
}
export const PageContentLayout: FC<PageContentLayout> = ({ children }) => {
  const [isSidePanelVisible, setSidePanelVisible] = useState(false);
  return (
    <div
      className={classNames(styles.pageContent, {
        [styles.open]: isSidePanelVisible,
      })}
    >
      <SidePanel
        onPanelToggle={setSidePanelVisible}
        isPanelOpen={isSidePanelVisible}
      />
      <div style={{ padding: '1.2rem', height: '100%', gridArea: 'content' }}>
        {children}
      </div>
    </div>
  );
};
