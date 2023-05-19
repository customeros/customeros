import React, { FC, useRef } from 'react';
import { useRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import { useMergeOrganizations } from '@spaces/hooks/useOrganization';
import EllipsesV from '@spaces/atoms/icons/EllipsesV';
import { Button } from '@spaces/atoms/button';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { OverlayPanel } from '@spaces/atoms/overlay-panel';
import styles from './finder-table.module.scss';
import { useMergeContacts } from '@spaces/hooks/useContact';
import { useRouter } from 'next/router';
import Check from '@spaces/atoms/icons/Check';

interface ActionColumnProps {
  onMerge: ({
    primaryId,
    mergeIds,
  }: {
    primaryId: string;
    mergeIds: Array<string>;
  }) => void;
  actions: Array<{
    label: string;
    command: () => void;
  }>;
}
export const ActionColumn: FC<ActionColumnProps> = ({ onMerge, actions }) => {
  const op = useRef(null);
  const [mode, setMode] = useRecoilState(tableMode);
  const [selectedItems, setSelectedItems] = useRecoilState(selectedItemsIds);

  const handleSave = async () => {
    const [primaryId, ...mergeIds] = selectedItems;
    return onMerge({ primaryId, mergeIds });
  };

  if (mode === 'MERGE') {
    if (selectedItems.length > 1) {
      return (
        <div className={styles.actionHeader}>
          <Button mode='primary' onClick={handleSave}>
            Merge
          </Button>
        </div>
      );
    }
    return (
      <div className={styles.actionHeader}>
        <IconButton
          mode='secondary'
          onClick={() => {
            setMode('PREVIEW');
            setSelectedItems([]);
          }}
          label='Done'
          icon={<Check height={24} />}
        ></IconButton>
      </div>
    );
  }

  return (
    <div className={styles.actionHeader}>
      <IconButton
        label='Actions'
        className={styles.actionsMenuButton}
        id={'finder-actions-dropdown-button'}
        mode='secondary'
        size='xxxxs'
        //@ts-expect-error revisit
        onClick={(e) => op?.current?.toggle(e)}
        icon={
          <EllipsesV
            height={24}
            width={24}
            style={{ transform: 'rotate(90deg)' }}
          />
        }
      />

      <OverlayPanel
        ref={op}
        style={{
          maxHeight: '400px',
          height: 'fit-content',
          overflowX: 'hidden',
          overflowY: 'auto',
          bottom: 0,
        }}
        model={actions}
      />
    </div>
  );
};
