import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './details-page-layout.module.scss';
import { ArrowLeft, IconButton, Ribbon } from '../../atoms';
import { useTenantName } from '../../../../hooks/useTenant';

interface DetailsPageLayout {
  children: ReactNode;
  onNavigateBack: () => void;
}
export const DetailsPageLayout: FC<DetailsPageLayout> = ({
  children,
  onNavigateBack,
}) => {
  const { data: tenant } = useTenantName();

  return (
    <div className={classNames(styles.layout)}>
      {tenant && (
        <Ribbon top={0}>
          When sending emails to your contacts, please BCC {tenant}
          @getopenline.com so that the email can be viewed in your Openline
          timeline.
        </Ribbon>
      )}

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
