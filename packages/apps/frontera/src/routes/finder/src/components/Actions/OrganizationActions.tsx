import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { useKeys, useKeyBindings } from 'rooks';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';
import { OrganizationStore } from '@store/Organizations/Organization.store';
import { ActionItem } from '@finder/components/Actions/components/ActionItem.tsx';

import { X } from '@ui/media/icons/X';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { OrganizationStage } from '@graphql/types';
import { Grid01 } from '@ui/media/icons/Grid01.tsx';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { useModKey } from '@shared/hooks/useModKey';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { TableInstance } from '@ui/presentation/Table';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';

interface TableActionsProps {
  selection: string[];
  focusedId?: string | null;
  isCommandMenuOpen: boolean;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<OrganizationStore>;
}

export const OrganizationTableActions = observer(
  ({
    table,
    enableKeyboardShortcuts,
    focusedId,
    selection,
    isCommandMenuOpen,
  }: TableActionsProps) => {
    const store = useStore();

    const [targetId, setTargetId] = useState<string | null>(null);

    const selectCount = selection?.length;

    const clearSelection = () => table.resetRowSelection();

    const handleMergeOrganizations = () => {
      store.ui.commandMenu.setType('MergeConfirmationModal');
      store.ui.commandMenu.setOpen(true);
    };

    const onHideOrganizations = () => {
      store.ui.commandMenu.setType('DeleteConfirmationModal');
      store.ui.commandMenu.setOpen(true);
      store.ui.commandMenu.setCallback(() => {
        clearSelection();
      });
    };

    const onCreateContact = () => {
      if (!focusedId) return;
      store.ui.commandMenu.setType('AddContactViaLinkedInUrl');

      store.ui.commandMenu.setOpen(true);
      store.ui.commandMenu.setContext({
        entity: 'Organization',
        ids: [focusedId],
      });
    };

    const onOpenCommandK = () => {
      if (selection?.length === 1) {
        store.ui.commandMenu.setType('OrganizationCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          ids: selection,
        });
      } else {
        store.ui.commandMenu.setType('OrganizationBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Organization',
          ids: selection,
        });
      }
      store.ui.commandMenu.setOpen(true);
    };

    useEffect(() => {
      if (selectCount === 1 && focusedId === selection[0]) {
        setTargetId(selection[0]);
      }

      if (selectCount < 1) {
        setTargetId(null);
        clearSelection();
      }
    }, [selectCount, focusedId]);

    const moveToAllOrgs = () => {
      if (!selectCount && !focusedId) return;

      if (!selectCount && focusedId) {
        store.organizations.updateStage(
          [focusedId],
          OrganizationStage.Unqualified,
        );

        return;
      }

      store.organizations.updateStage(selection, OrganizationStage.Unqualified);
      clearSelection();
    };

    const moveToTarget = () => {
      if (!selectCount && !focusedId) return;

      if (!selectCount && focusedId) {
        store.organizations.updateStage([focusedId], OrganizationStage.Target);

        return;
      }
      store.organizations.updateStage(selection, OrganizationStage.Target);
      clearSelection();
    };

    const moveToOpportunities = () => {
      if (!selectCount && !focusedId) return;

      if (!selectCount && focusedId) {
        store.organizations.updateStage([focusedId], OrganizationStage.Engaged);

        return;
      }
      store.organizations.updateStage(selection, OrganizationStage.Engaged);
      clearSelection();
    };

    const handleOpen = (type: CommandMenuType, property?: string) => {
      if (selection?.length > 1) {
        store.ui.commandMenu.setContext({
          ids: selection,
          entity: 'Organizations',
          property: property,
        });
      } else {
        store.ui.commandMenu.setContext({
          ids: [focusedId || ''],
          entity: 'Organization',
          property: property,
        });
      }

      if (selection?.length === 1) {
        store.ui.commandMenu.setContext({
          ids: selection,
          entity: 'Organization',
          property: property,
        });
      }
      store.ui.commandMenu.setType(type);
      store.ui.commandMenu.setOpen(true);
    };

    useKeyBindings(
      {
        u: moveToAllOrgs,
        t: moveToTarget,
        o: moveToOpportunities,
        c: (e) => {
          e.stopPropagation();
          e.preventDefault();

          if (selectCount > 1) return;
          onCreateContact();
        },
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
    useKeys(
      ['Shift', 'T'],
      (e) => {
        e.stopPropagation();
        e.preventDefault();
        handleOpen('ChangeTags');
      },
      { when: enableKeyboardShortcuts },
    );
    useKeys(
      ['Shift', 'R'],
      (e) => {
        e.stopPropagation();
        e.preventDefault();
        handleOpen('RenameOrganizationProperty', 'name');
      },
      { when: enableKeyboardShortcuts },
    );

    useModKey(
      'Backspace',
      () => {
        handleOpen('DeleteConfirmationModal');
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
          <ButtonGroup className='flex items-center translate-x-[-50%] justify-center bottom-[42px] *:border-none'>
            {selectCount && (
              <Tooltip
                className='p-1.5'
                label={
                  <div className='flex items-center text-sm'>Unselect all</div>
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
              onClick={onHideOrganizations}
              dataTest='org-actions-archive'
              icon={<Archive className='text-inherit size-3' />}
              tooltip={
                <div className='flex gap-1'>
                  <span className='text-sm'>Archive</span>
                  <div className='bg-gray-600  min-h-5 min-w-5 rounded flex justify-center items-center'>
                    {isUserPlatformMac() ? '⌘' : 'Ctrl'}
                  </div>
                  <div className='bg-gray-600  min-h-5 min-w-5 rounded flex justify-center items-center'>
                    <Delete className='text-inherit' />
                  </div>
                </div>
              }
            >
              Archive
            </ActionItem>
            {selectCount > 1 && (
              <ActionItem
                onClick={handleMergeOrganizations}
                tooltip={<span className='text-sm'>Merge</span>}
                icon={<Copy07 className='text-inherit size-3' />}
              >
                Merge
              </ActionItem>
            )}

            <ActionItem
              onClick={onOpenCommandK}
              dataTest='org-actions-commandk'
              icon={<Grid01 className='size-3 text-inherit' />}
            >
              Actions
            </ActionItem>
          </ButtonGroup>
        )}
      </>
    );
  },
);
