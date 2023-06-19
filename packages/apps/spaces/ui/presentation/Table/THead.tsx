import classNames from 'classnames';
import type { HeaderContext } from '@tanstack/react-table';

import Sort from '@spaces/atoms/icons/Sort';
import { IconButton } from '@spaces/atoms/icon-button';

import styles from './THead.module.scss';

interface THeadProps<T extends object> extends HeaderContext<T, unknown> {
  title: string;
  subTitle?: string;
  columnHasIcon?: boolean;
}

export const THead = <T extends object>({
  title,
  header,
  subTitle,
  columnHasIcon,
}: THeadProps<T>) => {
  const canSort = header.column.getCanSort();
  const isSorted = header.column.getIsSorted();
  const onToggleSort = header.column.getToggleSortingHandler();

  return (
    <div
      className={classNames(styles.thead, {
        [styles.withIcon]: columnHasIcon,
      })}
    >
      <div style={{ display: 'flex' }}>
        <span className={styles.title}>{title}</span>
        {canSort && (
          <IconButton
            isSquare
            mode='text'
            label='Sort'
            size='xxxxs'
            onClick={onToggleSort}
            icon={
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                <Sort
                  height={8}
                  color={isSorted === 'asc' ? '#3a3a3a' : '#969696'}
                  style={{
                    transform: 'rotate(180deg)',
                    marginBottom: 2,
                  }}
                />
                <Sort
                  height={8}
                  color={isSorted === 'desc' ? '#3a3a3a' : '#969696'}
                />
              </div>
            }
          />
        )}
      </div>
      {subTitle && <p className={styles.subTitle}>{subTitle}</p>}
    </div>
  );
};
