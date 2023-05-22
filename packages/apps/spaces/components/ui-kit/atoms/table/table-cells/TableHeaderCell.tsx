import styles from './table-cells.module.scss';
import { IconButton } from '@spaces/atoms/icon-button';
import Sort from '@spaces/atoms/icons/Sort';
import { FC, ReactNode } from 'react';
import { SortingDirection } from '../../../../../graphQL/__generated__/generated';
import classNames from 'classnames';

interface TableHeaderCellSortableProps {
  label: string;
  subLabel?: string;
  sortable: true;
  onSort: (order: 'ASC' | 'DESC') => void;
  order: 'ASC' | 'DESC';
}

interface TableHeaderCellNonSortableProps {
  label: string;
  subLabel?: string;
  sortable: false;
}

interface TableHeaderCellProps {
  label: string;
  subLabel?: string;
  hasAvatar?: boolean;
  children?: ReactNode;
}

export const TableHeaderCell: FC<TableHeaderCellProps> = ({
  label,
  subLabel,
  children,
  ...rest
}) => {
  return (
    <div
      className={classNames(styles.header, {
        [styles.labelWithAvatar]: rest.hasAvatar,
      })}
    >
      <div className={classNames(styles.label)}>
        {label}
        {children && children}
      </div>
      {subLabel && <span className={styles.subLabel}>{subLabel}</span>}
    </div>
  );
};
