import React from 'react';
import { TableHeaderCell } from '@spaces/atoms/table';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { useRecoilValue, useResetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
import { Times, TimesCircle } from '@spaces/atoms/icons';
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
      {mode === mergeMode && !!selectedIds.length && (
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
        <TableHeaderCell label={label} subLabel={subLabel} />
      </div>
    </div>
  );
};
