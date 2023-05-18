import React, { useRef } from 'react';
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

export const ActionColumn = ({ scope }: any) => {
  const router = useRouter();

  const op = useRef(null);
  const [mode, setMode] = useRecoilState(tableMode);
  const [selectedItems, setSelectedItems] = useRecoilState(selectedItemsIds);
  const { onMergeOrganizations } = useMergeOrganizations();
  const { onMergeContacts } = useMergeContacts();

  const handleSave = async () => {
    const [primaryId, ...mergeIds] = selectedItems;
    return mode === 'MERGE_CONTACT'
      ? onMergeContacts({
          primaryContactId: primaryId,
          mergedContactIds: mergeIds,
        })
      : onMergeOrganizations({
          primaryOrganizationId: primaryId,
          mergedOrganizationIds: mergeIds,
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
        Done
      </Button>
    );
  }

  const dropdownOptions = [];
  if (scope === 'MERGE_ORG') {
    dropdownOptions.push(
      {
        label: 'Add organization',
        command() {
          router.push('/organization/new');
        },
      },
      {
        label: 'Merge organizations',
        command() {
          return setMode('MERGE_ORG');
        },
      },
    );
  }
  if (scope === 'MERGE_CONTACT') {
    dropdownOptions.push(
      {
        label: 'Add contact',
        command() {
          router.push('/contact/new');
        },
      },
      {
        label: 'Merge contacts',
        command() {
          return setMode('MERGE_CONTACT');
        },
      },
    );
  }

  return (
    <div className={styles.actionHeader}>
      <IconButton
        label='Actions'
        className={styles.actionsMenuButton}
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
        model={dropdownOptions}
      />
    </div>
  );
};
