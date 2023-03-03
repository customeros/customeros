import React from 'react';
import { TableHeaderCell } from '../../ui-kit/atoms/table';
import { IconButton } from '../../ui-kit/atoms';
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

  return (
    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
      <TableHeaderCell label={label} subLabel={subLabel} />

      {mode === mergeMode && (
        <IconButton
          className={styles.deselectAllButton}
          mode='dangerLink'
          onClick={resetSelectedItems}
        />
      )}
    </div>
  );
};
