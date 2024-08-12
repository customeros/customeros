import { useState, useEffect } from 'react';

import { useKeys, useKeyBindings } from 'rooks';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import { X } from '@ui/media/icons/X';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { useModKey } from '@shared/hooks/useModKey';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { TableInstance } from '@ui/presentation/Table';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';
import { TableIdType, OrganizationStage } from '@graphql/types';
import { ActionItem } from '@organizations/components/Actions/ActionItem';

interface TableActionsProps {
  onHide: () => void;
  tableId?: TableIdType;
  focusedId?: string | null;
  onOpenCommandK: () => void;
  isCommandMenuOpen: boolean;
  onCreateContact: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<OrganizationStore>;
  handleOpen: (type: CommandMenuType) => void;
  onUpdateStage: (ids: string[], stage: OrganizationStage) => void;
}

export const OrganizationTableActions = ({
  table,
  onHide,
  tableId,
  onUpdateStage,
  enableKeyboardShortcuts,
  onCreateContact,
  focusedId,
  isCommandMenuOpen,
  onOpenCommandK,
  handleOpen,
}: TableActionsProps) => {
  const [targetId, setTargetId] = useState<string | null>(null);

  const selection = table.getState().rowSelection;

  const selectedIds = Object.keys(selection);

  const selectCount = selectedIds?.length;

  const clearSelection = () => table.resetRowSelection();

  const handleMergeOrganizations = () => {
    handleOpen('MergeConfirmationModal');
  };

  useEffect(() => {
    if (selectCount === 1 && focusedId === selectedIds[0]) {
      setTargetId(selectedIds[0]);
    }

    if (selectCount < 1) {
      setTargetId(null);
      clearSelection();
    }
  }, [selectCount, focusedId]);

  const moveToAllOrgs = () => {
    if (!selectCount && !focusedId) return;

    if (!selectCount && focusedId) {
      onUpdateStage([focusedId], OrganizationStage.Unqualified);

      return;
    }

    onUpdateStage(selectedIds, OrganizationStage.Unqualified);
    clearSelection();
  };

  const moveToTarget = () => {
    if (!selectCount && !focusedId) return;

    if (!selectCount && focusedId) {
      onUpdateStage([focusedId], OrganizationStage.Target);

      return;
    }
    onUpdateStage(selectedIds, OrganizationStage.Target);
    clearSelection();
  };

  const moveToOpportunities = () => {
    if (!selectCount && !focusedId) return;

    if (!selectCount && focusedId) {
      onUpdateStage([focusedId], OrganizationStage.Engaged);

      return;
    }
    onUpdateStage(selectedIds, OrganizationStage.Engaged);
    clearSelection();
  };

  const moveToLeads = (e: Event) => {
    if (!selectCount) return;
    e.preventDefault();
    e.stopPropagation();
    onUpdateStage(selectedIds, OrganizationStage.Lead);
    clearSelection();
  };

  useKeys(
    ['Shift', 'O'],
    () => {
      handleOpen('AssignOwner');
    },
    { when: enableKeyboardShortcuts },
  );

  useKeyBindings(
    {
      u: moveToAllOrgs,
      t: moveToTarget,
      o: moveToOpportunities,
      c: (e) => {
        e.stopPropagation();
        e.preventDefault();

        if (selectCount > 1) return;

        if (!targetId && focusedId) {
          setTargetId(focusedId);
        }
        onCreateContact();
      },
      l: (e) => tableId === TableIdType.Nurture && moveToLeads(e),
      Escape: clearSelection,
    },
    { when: enableKeyboardShortcuts },
  );
  useKeys(
    ['Shift', 'O'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('AssignOwner');
    },
    { when: enableKeyboardShortcuts },
  );

  useModKey(
    'Backspace',
    () => {
      onHide();
    },
    { when: enableKeyboardShortcuts },
  );

  useModKey(
    'v',
    () => {
      if (focusedId) {
        if (!targetId) {
          setTargetId(focusedId);
        }
        onCreateContact();
      }
    },
    {
      when:
        enableKeyboardShortcuts &&
        selectCount <= 1 &&
        typeof focusedId === 'string',
    },
  );

  return (
    <>
      {selectCount > 0 && !isCommandMenuOpen && (
        <ButtonGroup className='flex items-center translate-x-[-50%] justify-center bottom-[32px] *:border-none'>
          {selectCount && (
            <Tooltip
              className='p-1 pl-2'
              label={
                <div className='flex items-center text-sm'>
                  Open command menu
                  <div className='bg-gray-600 h-5 w-5 rounded-sm ml-3 mr-1 flex justify-center items-center'>
                    {isUserPlatformMac() ? '⌘' : 'Ctrl'}
                  </div>
                  <div className='bg-gray-600 text-xs h-5 w-5 rounded-sm flex justify-center items-center'>
                    K
                  </div>
                </div>
              }
            >
              <div className='bg-gray-700 px-3 py-2 rounded-s-lg'>
                <p
                  onClick={clearSelection}
                  className='flex text-gray-25 text-sm font-semibold text-nowrap leading-5 outline-dashed outline-1 rounded-[2px] outline-gray-400 pl-2 pr-1 hover:bg-gray-800 transition-colors cursor-pointer'
                >
                  {`${selectCount} selected`}
                  <span className='ml-1 inline-flex items-center'>
                    <X />
                  </span>
                </p>
              </div>
            </Tooltip>
          )}

          <ActionItem
            onClick={onHide}
            dataTest='org-actions-archive'
            icon={<Archive className='text-inherit size-3' />}
          >
            Archive
          </ActionItem>
          {selectCount > 1 && (
            <ActionItem
              onClick={handleMergeOrganizations}
              icon={<Copy07 className='text-inherit size-3' />}
            >
              Merge
            </ActionItem>
          )}
          <ActionItem
            onClick={onOpenCommandK}
            dataTest='org-actions-commandk'
            icon={
              <span className='text-inherit'>
                {isUserPlatformMac() ? '⌘' : 'Ctrl'}
              </span>
            }
          >
            Command
          </ActionItem>
        </ButtonGroup>
      )}
    </>
  );
};
