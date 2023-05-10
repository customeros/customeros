import React from 'react';
import { Skeleton } from '@spaces/atoms/skeleton';
import styles from '../contact-details.module.scss';

export const ContactDetailsSkeleton: React.FC = () => {
  return (
    <div className={styles.contactDetails} style={{ width: '100%' }}>
      <div className={styles.header} style={{ width: '100%' }}>
        <div className={styles.photo} style={{ background: '#ddd' }}>
          <div style={{ width: '40px', height: '40px' }} />
        </div>
        <div className={styles.name} style={{ width: '80%' }}>
          <div>
            <Skeleton height='20px' />
          </div>

          <div className={styles.jobRole}>
            <Skeleton />
          </div>

          {
            <div
              className={styles.source}
              style={{
                width: '100px',
              }}
            >
              <Skeleton height='8px' />
            </div>
          }
        </div>
      </div>
    </div>
  );
};
