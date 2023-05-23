import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import styles from './page-content-layout.module.scss';
import { SidePanel } from '@spaces/organisms/side-panel';

interface PageContentLayout {
  children: ReactNode;
}

export const PageContentLayout: FC<PageContentLayout> = ({ children }) => {
  return (
    <div className={styles.pageContent}>
      <SidePanel />
      <div style={{ padding: '1.2rem', height: '100%', gridArea: 'content', overflowX: 'hidden', overflowY: 'auto' }}>
        {children}
      </div>
    </div>
  );
};
