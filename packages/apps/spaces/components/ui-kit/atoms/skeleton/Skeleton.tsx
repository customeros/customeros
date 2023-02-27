import React from 'react';
import styles from './skeleton.module.scss';

export const Skeleton = ({
  height = 'auto',
  width = '80%',
}: {
  height?: string;
  width?: string;
}): JSX.Element => {
  return (
    <span
      className={styles.skeleton}
      style={{ height: height, width, minHeight: '8px' }}
    />
  );
};
