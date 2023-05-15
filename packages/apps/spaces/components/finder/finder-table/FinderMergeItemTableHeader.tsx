import React from 'react';
import { TableHeaderCell } from '@spaces/atoms/table';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { useRecoilValue, useResetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
export const FinderMergeItemTableHeader: React.FC<{
  mergeMode: 'MERGE_ORG' | 'MERGE_CONTACT';
  label: string;
  subLabel: string;
}> = ({ mergeMode, label, subLabel }) => {
  const resetSelectedItems = useResetRecoilState(selectedItemsIds);
  const mode = useRecoilValue(tableMode);
  const selectedIds = useRecoilValue(selectedItemsIds);

  return (
    <div style={{ display: 'flex' }}>
      <div className={styles.checkboxContainer}>
        {mode === mergeMode && !!selectedIds.length && (
          <IconButton
            size='xxxs'
            label='Deselect all items'
            className={styles.deselectAllButton}
            mode='dangerLink'
            onClick={resetSelectedItems}
          />
        )}
      </div>
      <div className={styles.finderCell}>
        <TableHeaderCell label={label} subLabel={subLabel} />
      </div>
    </div>
  );
};
