import React from 'react';
import { Skeleton } from '../../../ui-kit/atoms/skeleton';
import styles from '../organization-timeline.module.scss';
import { Envelope, Phone } from '../../../ui-kit';

export const OrganizationTimelineSkeleton: React.FC = () => {
  const rows = Array(2)
    .fill('')
    .map((e, i) => i + 1);
  return (
    <div>
      {rows.map((row, id) => (
        <div
          key={`organization-contacts-skeleton-${row}-${id}`}
          className={styles.contactItem}
        >
          <div className={styles.contactDetails}>
            <div className={styles.name}>
              <Skeleton height={'16px'} />
            </div>
            <div
              className={styles.detailsContainer}
              style={{ maxWidth: '50%' }}
            >
              <Envelope className={styles.icon} />
              <Skeleton height={'12px'} />
            </div>
            <div
              className={styles.detailsContainer}
              style={{ maxWidth: '50%' }}
            >
              <Phone className={styles.icon} />
              <Skeleton height={'12px'} />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
