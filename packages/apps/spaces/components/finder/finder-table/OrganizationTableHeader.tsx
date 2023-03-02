import React from 'react';
import { TableHeaderCell } from '../../ui-kit/atoms/table';
import { IconButton, Times } from '../../ui-kit/atoms';
import { useRecoilValue, useResetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import styles from './finder-table.module.scss';
export const OrganizationTableHeader: React.FC = () => {
  const resetSelectedItems = useResetRecoilState(selectedItemsIds);
  const mode = useRecoilValue(tableMode);

  return (
    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
      <TableHeaderCell label='Organization' subLabel='Industry' />

      {mode === 'MERGE_ORG' && (
        <IconButton
          className={styles.deselectAllButton}
          mode='dangerLink'
          onClick={resetSelectedItems}
        />
      )}
    </div>
  );
};
