import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './details-page-layout.module.scss';

interface DetailsPageLayout {
  children: ReactNode;
}
export const DetailsPageLayout: FC<DetailsPageLayout> = ({ children }) => {
  return <div className={classNames(styles.layout)}>{children}</div>;
};
