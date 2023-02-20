import styles from './dashboard-table-header-label.module.scss';

export const DashboardTableHeaderLabel = ({
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
