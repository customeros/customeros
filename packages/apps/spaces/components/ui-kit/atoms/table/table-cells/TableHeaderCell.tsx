import styles from './table-cells.module.scss';
import { FC, ReactNode } from 'react';
import classNames from 'classnames';

interface TableHeaderCellProps {
  label: string;
  withIcon?: boolean;
  subLabel?: string;
  children?: ReactNode;
}

export const TableHeaderCell: FC<TableHeaderCellProps> = ({
  label,
  subLabel,
  children,
  withIcon,
}) => {
  return (
    <div
      className={classNames(styles.header, {
        [styles.labelWithIcon]: withIcon,
      })}
    >
      <div className={classNames(styles.label, styles.headerLabel)}>
        {label}
        {children && children}
      </div>
      {subLabel && <span className={styles.subLabel}>{subLabel}</span>}
    </div>
  );
};
