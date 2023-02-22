import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './page-content-layout.module.scss';

interface PageContentLayout {
  isPanelOpen: boolean;
  isSideBarShown: boolean;
  children: ReactNode;
}
export const PageContentLayout: FC<PageContentLayout> = ({
  children,
  isPanelOpen,
}) => {
  return (
    <div
      className={classNames(styles.pageContent, { [styles.open]: isPanelOpen })}
    >
      {children}
    </div>
  );
};
