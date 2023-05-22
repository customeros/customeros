import React from 'react';
import styles from '../table.module.scss';
import { Skeleton } from '../../skeleton';
import { Column } from '../types';
import classNames from 'classnames';
import { TableHeaderCell } from '@spaces/atoms/table';
import { TableContentSkeleton } from '@spaces/atoms/table/skeletons/TableContentSkeleton';
interface TableSkeletonProps {
  columns: number;
}

export const TableSkeleton = ({ columns }: TableSkeletonProps): JSX.Element => {
  const columnsNew = Array(4)
    .fill('')
    .map((e, i) => ({ label: i + 1, width: `${100 / columns}%` }));
  return (
    <>
      <div className={styles.itemCounter}>
        <span>Total items:</span>
        <Skeleton height='16px' width='16px' isSquare />
      </div>

      <table className={styles.table}>
        <thead className={styles.header}>
          <tr>
            {columnsNew?.map(({ label, width }) => {
              return (
                <th
                  key={`table-header-${label}`}
                  style={{
                    width: width || 'auto',
                    minWidth: width || 'auto',
                    maxWidth: width || 'auto',
                  }}
                >
                  <Skeleton height='8px' width='80px' />
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody className={styles.body}>
          <TableContentSkeleton columns={columnsNew} />
        </tbody>
      </table>
    </>
  );
};
