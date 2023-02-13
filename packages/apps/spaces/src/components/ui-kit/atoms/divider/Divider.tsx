import React from 'react';
import styles from './divider.module.css';

export const Divider: React.FC<Partial<HTMLDivElement>> = ({ ...props }) => {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  return <div {...props} className={styles.divider} />;
};
