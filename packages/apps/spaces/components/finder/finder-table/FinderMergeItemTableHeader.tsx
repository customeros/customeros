import React from 'react';
import { TableHeaderCell } from '@spaces/atoms/table';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { useRecoilValue, useResetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import Times from '@spaces/atoms/icons/Times';

export const FinderMergeItemTableHeader: React.FC<{
  label: string;
  subLabel: string;
  withIcon?: boolean;
  children: React.ReactNode;
}> = ({ label, subLabel, children, withIcon }) => {
  const resetSelectedItems = useResetRecoilState(selectedItemsIds);
  const mode = useRecoilValue(tableMode);
  const selectedIds = useRecoilValue(selectedItemsIds);

  return (
    <div style={{ display: 'flex' }}>
      {(mode === 'MERGE' || mode === 'ARCHIVE') && !!selectedIds.length && (
        <IconButton
          size='xxxs'
          label='Deselect all items'
          className={styles.deselectAllButton}
          mode='dangerLink'
          icon={<Times height={14} width={14} />}
          onClick={resetSelectedItems}
        />
      )}
      <div className={styles.finderCell}>
        <TableHeaderCell label={label} subLabel={subLabel} withIcon={withIcon}>
          {children}
        </TableHeaderCell>
      </div>
    </div>
  );
};
