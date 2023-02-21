import React, { useEffect } from 'react';
import { DashboardTableHeaderLabel } from './dashboard-table-header-label';
import { useVirtual } from 'react-virtual';
import styles from './table.module.scss';
import { Skeleton } from '../skeleton';
import { Column } from './types';
import { TableSkeleton } from './TableSkeleton';

interface TableProps<T> {
  data: Array<T> | null;
  onFetchNextPage: () => void;
  isFetching: boolean;
  columns: Array<Column>;
  totalItems: number;
}

export const Table = <T,>({
  columns,
  data,
  totalItems,
  isFetching,
  onFetchNextPage,
}: TableProps<T>) => {
  const parentRef = React.useRef(null);
  const rowVirtualizer = useVirtual({
    size: totalItems,
    parentRef,
    estimateSize: React.useCallback(() => 70, []),
    overscan: 5,
  });

  useEffect(() => {
    const [lastItem] = [...rowVirtualizer.virtualItems].reverse();
    if (!lastItem || !data) {
      return;
    }

    if (lastItem.index >= data?.length - 1 && !isFetching) {
      onFetchNextPage();
    }
  }, [
    totalItems,
    onFetchNextPage,
    isFetching,
    rowVirtualizer.virtualItems,
    data,
  ]);
  return (
    <table className={styles.table}>
      <thead className={styles.header}>
        <tr>
          {columns?.map(({ label, subLabel, width }) => {
            return (
              <th
                key={`header-${label}`}
                style={{ width }}
                data-th={label}
                data-th2={subLabel}
              >
                <DashboardTableHeaderLabel
                  label={label}
                  subLabel={subLabel || ''}
                />
              </th>
            );
          })}
        </tr>
      </thead>
      <tbody ref={parentRef} className={styles.body}>
        {(!data || !data.length) && <TableSkeleton columns={columns} />}
        {/* SHOW TABLE*/}
        {!!data &&
          rowVirtualizer.virtualItems.map((virtualRow) => {
            const element = data[virtualRow.index];

            return (
              <tr
                key={virtualRow.key}
                data-index={virtualRow.index}
                ref={virtualRow.measureRef}
                className={styles.row}
                style={{
                  minHeight: `${virtualRow.size}px`,
                  transform: `translateY(${virtualRow.start}px)`,
                }}
              >
                {columns.map(({ template, width, label }) => (
                  <td
                    key={`table-row-${label}`}
                    style={{
                      width: width || 'auto',
                      maxWidth: width || 'auto',
                    }}
                  >
                    {element && template(element)}
                    {!element && <Skeleton />}
                  </td>
                ))}
              </tr>
            );
          })}
      </tbody>
    </table>
  );
};
