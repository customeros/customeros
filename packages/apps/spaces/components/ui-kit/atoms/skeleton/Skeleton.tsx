import React from 'react';
import styles from './skeleton.module.scss';
import classNames from 'classnames';

export const Skeleton = ({
  height = 'auto',
  width = '80%',
  className,
}: {
  height?: string;
  width?: string;
  className?: string;
}): JSX.Element => {
  return (
    <span
      className={classNames(styles.skeleton, className)}
      style={{ height: height, width, minHeight: height || '8px' }}
    />
  );
};
