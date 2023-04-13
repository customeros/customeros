import React from 'react';
import { ListSkeleton } from './ListSkeleton';
import styles from '../contact-communication-details.module.scss';

export const ContactCommunicationDetailsSkeleton: React.FC = () => {
  return (
    <div className={styles.contactDetails}>
      <ListSkeleton id='email-list' />
      <div className={styles.divider} />
      <ListSkeleton id='numbers-list' />
    </div>
  );
};
