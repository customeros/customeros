import React from 'react';
import styles from './inline-loader.module.scss';

export const InlineLoader: React.FC<{ label?: string }> = ({
  label = 'Saving...',
}) => {
  return (
    <div aria-label={label} className={styles.dot_flashing_container}>
      <div className={styles.dot_flashing} />
      <div className={styles.dot_flashing} />
      <div className={styles.dot_flashing} />
    </div>
  );
};
