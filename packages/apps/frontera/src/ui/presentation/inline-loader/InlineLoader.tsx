import React from 'react';

import styles from './inline-loader.module.scss';

export const InlineLoader: React.FC<{ label?: string; color?: string }> = ({
  label = 'Saving...',
  color = '#9880ff',
}) => {
  return (
    <div
      title={label}
      aria-label={label}
      className={styles.dot_flashing_container}
      // @ts-expect-error fixme
      style={{ '--flashing-dot-color': color }}
    >
      <div className={styles.dot_flashing} />
      <div className={styles.dot_flashing} />
      <div className={styles.dot_flashing} />
    </div>
  );
};
