import React from 'react';
import styles from './skeleton.module.scss';
import classNames from 'classnames';

export const Skeleton = ({
  height = '8px',
  width = '100%',
  className,
  isSquare,
}: {
  height?: string;
  width?: string;
  className?: string;
  isSquare?: boolean;
}): JSX.Element => {
  return (
    <span
      className={classNames(styles.skeleton, className, {
        [styles.squareSkeleton]: isSquare,
      })}
      style={{ height: height, width, minHeight: height || '8px' }}
    />
  );
};
