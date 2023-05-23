import React from 'react';
import styles from '../table.module.scss';
import { Skeleton } from '../../skeleton';
import { Column } from '../types';
import classNames from 'classnames';

interface TableSkeletonProps {
  columns: Array<
    Column | { width: string; label: string | number; isLast?: boolean }
  >;
}

export const TableContentSkeleton = ({
  columns,
}: TableSkeletonProps): JSX.Element => {
  const rows = Array(4)
    .fill('')
    .map((e, i) => ({ label: i + 1, width: '20%' }));
  return (
    <>
      {rows.map((n) => (
        <tr
          key={`skeleton-row-${n.label}`}
          className={classNames(styles.row, styles.staticRow)}
        >
          {columns.map(({ width, label, ...rest }, index) => (
            <td
              key={`table-skeleton-${label}-${index}`}
              className={classNames({
                [styles.actionCell]: rest?.isLast,
              })}
              style={{
                width: width || 'auto',
                maxWidth: width || 'auto',
              }}
            >
              <Skeleton />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
};
