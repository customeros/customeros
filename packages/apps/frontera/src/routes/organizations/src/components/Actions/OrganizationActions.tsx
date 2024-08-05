import { useState, useEffect } from 'react';

import { useKeyBindings } from 'rooks';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import { X } from '@ui/media/icons/X';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { Archive } from '@ui/media/icons/Archive';
import { UserX01 } from '@ui/media/icons/UserX01';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { useModKey } from '@shared/hooks/useModKey';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { TableInstance } from '@ui/presentation/Table';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { TableIdType, OrganizationStage } from '@graphql/types';
import { LinkedInSolid } from '@ui/media/icons/LinkedInSolid.tsx';
import { ActionItem } from '@organizations/components/Actions/ActionItem.tsx';
import { EditTagsModal } from '@organizations/components/Actions/components/EditTagsModal.tsx';

interface TableActionsProps {
  onHide: () => void;
  tableId?: TableIdType;
  focusedId?: string | null;
  isCommandMenuOpen: boolean;
  onCreateContact: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<OrganizationStore>;
  onMerge: (primaryId: string, mergeIds: string[]) => void;
  onUpdateStage: (ids: string[], stage: OrganizationStage) => void;
}

export const OrganizationTableActions = ({
  table,
  onHide,
  onMerge,
  tableId,
  onUpdateStage,
  enableKeyboardShortcuts,
  onCreateContact,
  focusedId,
  isCommandMenuOpen,
}: TableActionsProps) => {
  const {
    open: isTagEditOpen,
    onOpen: onOpenTagEdit,
    onClose: onCloseTagEdit,
  } = useDisclosure({ id: 'organization-tag-actions' });
  const [targetId, setTargetId] = useState<string | null>(null);

  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection);
  const selectCount = selectedIds.length;
  const clearSelection = () => table.resetRowSelection();

  const handleMergeOrganizations = () => {
    const mergeIds = selectedIds.filter((id) => id !== targetId);

    if (!targetId || !mergeIds.length) return;

    onMerge(targetId, mergeIds);
    clearSelection();
  };

  useEffect(() => {
    if (selectCount === 1) {
      setTargetId(selectedIds[0]);
    }

    if (selectCount < 1) {
      setTargetId(null);
    }
  }, [selectCount]);

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
        tableId === TableIdType.Nurture && onCreateContact();
      },
      l: (e) => tableId === TableIdType.Nurture && moveToLeads(e),
      Escape: clearSelection,
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
        tableId === TableIdType.Nurture &&
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
          {tableId &&
            [
              TableIdType.Leads,
              TableIdType.Nurture,
              TableIdType.Organizations,
            ].includes(tableId) && (
              <ActionItem
                shortcutKey='U'
                onClick={moveToAllOrgs}
                icon={<UserX01 className='text-inherit size-3' />}
                tooltip={'Change to Unqualified and move to All orgs'}
              >
                Unqualify
              </ActionItem>
            )}

          {tableId &&
            [TableIdType.Leads, TableIdType.Organizations].includes(
              tableId,
            ) && (
              <ActionItem
                shortcutKey='T'
                onClick={moveToTarget}
                tooltip='Change to Target and move to Targets'
                icon={<HeartHand className='text-inherit size-3' />}
              >
                Target
              </ActionItem>
            )}
          {tableId &&
            selectedIds.length === 1 &&
            [TableIdType.Nurture, TableIdType.Organizations].includes(
              tableId,
            ) && (
              <ActionItem
                shortcutKey='C'
                onClick={onCreateContact}
                tooltip='Add contact via LinkedIn URL'
                icon={
                  <LinkedInSolid className='text-inherit size-4 text-inherit ' />
                }
              >
                Add LinkedIn contact
              </ActionItem>
            )}

          {tableId &&
            [TableIdType.Leads, TableIdType.Nurture].includes(tableId) && (
              <ActionItem
                onClick={onOpenTagEdit}
                icon={<Tag01 className='text-inherit size-4 text-inherit ' />}
              >
                Edit tags
              </ActionItem>
            )}
          {tableId &&
            [
              TableIdType.Leads,
              TableIdType.Nurture,
              TableIdType.Organizations,
            ].includes(tableId) && (
              <ActionItem
                shortcutKey='O'
                onClick={moveToOpportunities}
                tooltip='Change to Engaged and move to Opportunities'
                icon={<CoinsStacked01 className='text-inherit size-3' />}
              >
                Opportunity
              </ActionItem>
            )}
        </ButtonGroup>
      )}

      <EditTagsModal
        isOpen={isTagEditOpen}
        onClose={onCloseTagEdit}
        selectedIds={selectedIds}
        clearSelection={clearSelection}
      />
    </>
  );
};
