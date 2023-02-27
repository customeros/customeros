import styles from './table-cells.module.scss';

export const TableHeaderCell = ({
  label,
  subLabel,
}: {
  label: string;
  subLabel: string;
}) => {
  return (
    <div className={styles.header}>
      <span>{label}</span>
      {subLabel && <span className={styles.subLabel}>{subLabel}</span>}
    </div>
  );
};
