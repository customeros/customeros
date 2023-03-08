import React from 'react';
import { TableHeaderCell } from '../../ui-kit/atoms/table';
import { Button, IconButton } from '../../ui-kit/atoms';
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
