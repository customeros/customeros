import React from 'react';
import styles from './skeleton.module.scss';

export const Skeleton = ({
  height = 'auto',
}: {
  height?: string;
}): JSX.Element => {
  return (
    <span
      className={styles.skeleton}
      style={{ height: height, minHeight: '8px' }}
    />
  );
};
