import styles from './inline-loader.module.scss';

export const InlineLoader = ({
  label = 'Saving...',
  color = '#9880ff',
}: {
  label?: string;
  color?: string;
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
