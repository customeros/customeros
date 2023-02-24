import React from 'react';
import { Skeleton } from '../../../ui-kit/atoms/skeleton';
import { ListSkeleton } from './ListSkeleton';
import styles from '../contact-communication-details.module.scss';

export const ContactCommunicationDetailsSkeleton: React.FC = () => {
  return (
    <div className={styles.contactDetails}>
      <div className={styles.buttonWrapper}>
        <Skeleton height={'30px'} />
      </div>
      <ListSkeleton id='email-list' />
      <div className={styles.divider} />
      <ListSkeleton id='numbers-list' />
    </div>
  );
};
