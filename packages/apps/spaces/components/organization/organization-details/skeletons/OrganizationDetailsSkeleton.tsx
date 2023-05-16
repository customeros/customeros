import React from 'react';
import { Skeleton } from '@spaces/atoms/skeleton';
import styles from '../organization-details.module.scss';
import { OrganizationCommunicationDetailsSkeleton } from './OrganizationCommunicationDetailsSkeleton';
import { OrganizationSubsidiariesSkeleton } from '@spaces/organization/organization-details/subsidiaries/skeletons';

export const OrganizationDetailsSkeleton: React.FC = () => {
  return (
    <div>
      <div className={styles.organizationDetails}>
        <div className={styles.bg}>
          <div style={{ maxWidth: '80%' }}>
            <h1 className={styles.name}>
              <Skeleton height={'20px'} />
            </h1>
            <span className={styles.industry} style={{ maxWidth: '50%' }}>
              <Skeleton />
            </span>
          </div>

          <p className={styles.description}>
            <Skeleton />
            <Skeleton />
          </p>
        </div>
      </div>
      <OrganizationCommunicationDetailsSkeleton />
      <OrganizationSubsidiariesSkeleton />
    </div>
  );
};
