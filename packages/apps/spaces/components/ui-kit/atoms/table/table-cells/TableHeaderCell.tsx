import styles from './table-cells.module.scss';
import { FC, ReactNode } from 'react';
import classNames from 'classnames';

interface TableHeaderCellProps {
  label: string;
  subLabel?: string;
  children?: ReactNode;
}

export const TableHeaderCell: FC<TableHeaderCellProps> = ({
  label,
  subLabel,
  children,
}) => {
  return (
    <div
      className={classNames(styles.header)}
    >
      <div className={classNames(styles.label)}>
        {label}
        {children && children}
      </div>
      {subLabel && <span className={styles.subLabel}>{subLabel}</span>}
    </div>
  );
};
