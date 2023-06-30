import '@openline-ai/openline-web-chat/dist/esm/index.css';
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './details-page-layout.module.scss';
import { Ribbon } from '@spaces/atoms/ribbon';
import { useTenantName } from '@spaces/hooks/useTenant';
import { useRecoilValue } from 'recoil';
import { tenantName } from '../../../../state/userData';

interface DetailsPageLayout {
  children: ReactNode;
}

export const DetailsPageLayout: FC<DetailsPageLayout> = ({ children }) => {
  useTenantName();
  const tenant = useRecoilValue(tenantName);

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
