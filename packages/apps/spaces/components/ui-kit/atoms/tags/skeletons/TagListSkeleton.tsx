import React from 'react';
import styles from '../tags.module.scss';
import { Skeleton } from '../../skeleton';

export const TagListSkeleton: React.FC = () => {
  return (
    <div className={styles.tagsList}>
      <Skeleton height={'21px'} width={'60px'} />
    </div>
  );
};
