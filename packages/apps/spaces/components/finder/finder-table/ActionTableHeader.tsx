import React, { useRef } from 'react';
import { useRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../state';
import { useMergeOrganizations } from '../../../hooks/useOrganization';
import { Button, EllipsesV, Tooltip } from '../../ui-kit';
import { IconButton } from '../../ui-kit/atoms';
import { OverlayPanel } from '../../ui-kit/atoms/overlay-panel';
import styles from './finder-table.module.scss';

export const ActionColumn = () => {
  const op = useRef(null);
  const [mode, setMode] = useRecoilState(tableMode);
  const [selectedItems, setSelectedItems] = useRecoilState(selectedItemsIds);
  const { onMergeOrganizations } = useMergeOrganizations();

  const handleSave = async () => {
    const [primaryOrganizationId, ...mergedOrganizationIds] = selectedItems;
    await onMergeOrganizations({
      primaryOrganizationId,
      mergedOrganizationIds,
    });
  };

  if (mode === 'MERGE_ORG' || mode === 'MERGE_CONTACT') {
    if (selectedItems.length > 1) {
      return (
        <Button mode='primary' onClick={handleSave}>
          Merge
        </Button>
      );
    }
    return (
      <Button
        mode='secondary'
        onClick={() => {
          setMode('PREVIEW');
          setSelectedItems([]);
        }}
      >
        Cancel
      </Button>
    );
  }

  return (
    <div className={styles.actionHeader}>
      <Tooltip
        content='Actions'
        target='#finder-actions-dropdown-button'
        position='top'
        showDelay={300}
        autoHide={false}
      />
      <IconButton
        id={'finder-actions-dropdown-button'}
        mode='secondary'
        size={'xxxs'}
        //@ts-expect-error revisit
        onClick={(e) => op?.current?.toggle(e)}
        icon={<EllipsesV style={{ transform: 'rotate(90deg)' }} />}
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
        model={[
          {
            label: 'Merge organizations',
            command() {
              return setMode('MERGE_ORG');
            },
          },
          {
            label: 'Merge contacts',
            disabled: true,
            command() {
              return setMode('MERGE_CONTACT');
            },
          },
        ]}
      />
    </div>
  );
};
