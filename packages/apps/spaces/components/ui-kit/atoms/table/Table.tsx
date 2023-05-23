import React, { useEffect } from 'react';
import { TableHeaderCell } from './table-cells';
import { useVirtual } from 'react-virtual';
import styles from './table.module.scss';
import { Skeleton } from '../skeleton';
import { Column } from './types';
import { TableContentSkeleton } from './skeletons/TableContentSkeleton';
import { SearchMinus } from '../icons';
import classNames from 'classnames';

interface TableProps<T> {
  data: Array<T> | null;
  onFetchNextPage: () => void;
  isFetching: boolean;
  columns: Array<Column>;
  totalItems: number;
}

export const Table = <T,>({
  columns,
  data = [],
  totalItems,
  isFetching,
  onFetchNextPage,
}: TableProps<T>) => {
  const parentRef = React.useRef(null);
  const rowVirtualizer = useVirtual({
    size: totalItems,
    parentRef,
    estimateSize: React.useCallback(() => 54, []),
    overscan: 5,
  });
  useEffect(() => {
    const [lastItem] = [...rowVirtualizer.virtualItems].reverse();
    if (!lastItem || !data) {
      return;
    }

    if (
      lastItem.index >= data?.length - 1 &&
      data.length < totalItems &&
      !isFetching
    ) {
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
    <>
      <div className={styles.itemCounter}>
        <span>Total items:</span>
        {totalItems}
      </div>

      <table className={styles.table}>
        <thead className={styles.header}>
          <tr>
            {columns?.map(({ label, id, subLabel, width, ...rest }) => {
              if (typeof label !== 'string') {
                return (
                  <th
                    key={`header-${id}`}
                    style={{
                      width: width || 'auto',
                      minWidth: width || 'auto',
                      maxWidth: width || 'auto',
                    }}
                    className={classNames({
                      [styles.actionCell]: rest?.isLast,
                    })}
                    data-th={label}
                    data-th2={subLabel}
                  >
                    {label}
                  </th>
                );
              }

              return (
                <th
                  key={`header-${id}`}
                  style={{
                    width: width || 'auto',
                    minWidth: width || 'auto',
                    maxWidth: width || 'auto',
                  }}
                  data-th={label}
                  data-th2={subLabel}
                  className={classNames({
                    [styles.actionCell]: rest?.isLast,
                  })}
                >
                  <TableHeaderCell label={label} subLabel={subLabel || ''} />
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody ref={parentRef} className={styles.body}>
          {/* SHOW TABLE*/}
          {!totalItems && !isFetching && (
            <tr className={styles.noResultsInfo}>
              <td>
                <SearchMinus />
                <span>No results</span>
              </td>
            </tr>
          )}
          {rowVirtualizer.virtualItems.map((virtualRow) => {
            const element = data?.[virtualRow.index];
            return (
              <tr
                key={virtualRow.key}
                data-index={virtualRow.index}
                ref={virtualRow.measureRef}
                className={classNames(styles.row, {
                  [styles.odd]: virtualRow.index % 2 !== 0,
                })}
                style={{
                  // padding: `5px 0px`,
                  minHeight: `${virtualRow.size}px`,
                  transform: `translateY(${virtualRow.start}px)`,
                }}
              >
                {columns.map(({ template, width, id, ...rest }) => (
                  <td
                    key={`table-row-${id}`}
                    className={classNames({
                      [styles.actionCell]: rest?.isLast,
                    })}
                    style={{
                      width: width || 'auto',
                      minWidth: width || 'auto',
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
          {isFetching && !totalItems && (
            <TableContentSkeleton columns={columns} />
          )}
        </tbody>
      </table>
    </>
  );
};
