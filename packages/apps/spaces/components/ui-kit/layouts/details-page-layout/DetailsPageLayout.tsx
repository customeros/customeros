import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './details-page-layout.module.scss';
import { ArrowLeft, Button, IconButton } from '../../atoms';

interface DetailsPageLayout {
  children: ReactNode;
  onNavigateBack: () => void;
}
export const DetailsPageLayout: FC<DetailsPageLayout> = ({
  children,
  onNavigateBack,
}) => {
  return (
    <div className={classNames(styles.layout)}>
      <div className={styles.backButton}>
        <IconButton
          mode='secondary'
          icon={<ArrowLeft />}
          onClick={onNavigateBack}
        />
      </div>

      {children}
    </div>
  );
};
