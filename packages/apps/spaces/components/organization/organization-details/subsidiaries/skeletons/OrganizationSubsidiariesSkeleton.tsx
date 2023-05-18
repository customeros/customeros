import React from 'react';
import { Skeleton } from '@spaces/atoms/skeleton';
import styles from '../organization-subsidiaries.module.scss';
import classNames from 'classnames';

export const OrganizationSubsidiariesSkeleton: React.FC = () => {
  return (
    <article className={styles.subsidiary_section}>
      <h1 className={styles.subsidiary_header}>Branches</h1>

      <Skeleton
        className={classNames(styles.subsidiary, styles.subsidiary_skeleton)}
        height={'8px'}
        width={'140px'}
      />
      <Skeleton
        className={classNames(styles.subsidiary, styles.subsidiary_skeleton)}
        height={'8px'}
        width={'140px'}
      />
      <Skeleton
        className={classNames(styles.subsidiary, styles.subsidiary_skeleton)}
        height={'8px'}
        width={'140px'}
      />
      <Skeleton
        className={classNames(styles.subsidiary, styles.subsidiary_skeleton)}
        height={'8px'}
        width={'140px'}
      />
    </article>
  );
};
