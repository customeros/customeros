import styles from './table-cells.module.scss';
import { IconButton } from '@spaces/atoms/icon-button';
import Sort from '@spaces/atoms/icons/Sort';
import { FC } from 'react';
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

interface TableHeaderCellProps<T extends boolean> {
  label: string;
  subLabel?: string;
  sortable: T;
  onSort?: T extends true ? (direction: SortingDirection) => void : undefined;
  direction?: T extends true ? SortingDirection : undefined;
}

export const TableHeaderCell: FC<any> = ({
  label,
  subLabel,
  sortable,
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
        {sortable && (
          <IconButton
            isSquare
            mode='text'
            onClick={() => {
              rest.onSort(
                rest.direction === SortingDirection.Asc
                  ? SortingDirection.Desc
                  : SortingDirection.Asc,
              );
            }}
            label='SORT'
            size={'xxxxs'}
            icon={
              <Sort
                height={10}
                color='#969696'
                style={{
                  transform:
                    rest.direction === SortingDirection.Asc
                      ? 'rotate(180deg)'
                      : '',
                }}
              />
            }
          />
        )}
      </div>
      {subLabel && <span className={styles.subLabel}>{subLabel}</span>}
    </div>
  );
};
