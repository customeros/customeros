import React from 'react';
import { Skeleton } from '../../../ui-kit/atoms/skeleton';
import styles from '../contact-communication-details.module.scss';

export const ListSkeleton = ({ id }: { id: string }) => {
  const rows = Array(2)
    .fill('')
    .map((e, i) => i + 1);
  return (
    <div>
      {rows.map((row) => (
        <div key={`${row}-${id}`} className={styles.communicationItem}>
          <div className={styles.detailsList}>
            <Skeleton height='12px' width='60px' />
          </div>
          <div className={styles.detailsList} style={{ marginLeft: '8px' }}>
            <Skeleton height='12px' width={'200px'} />
          </div>
        </div>
      ))}
    </div>
  );
};
