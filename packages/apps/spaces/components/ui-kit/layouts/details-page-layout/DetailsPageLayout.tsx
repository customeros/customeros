import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './details-page-layout.module.scss';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { Ribbon } from '@spaces/atoms/ribbon';
import ArrowLeft from '@spaces/atoms/icons/ArrowLeft';
import { useTenantName } from '@spaces/hooks/useTenant';

interface DetailsPageLayout {
  children: ReactNode;
}

export const DetailsPageLayout: FC<DetailsPageLayout> = ({ children }) => {
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
      {children}
    </div>
  );
};
