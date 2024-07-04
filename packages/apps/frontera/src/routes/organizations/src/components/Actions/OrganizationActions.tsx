import type { Store } from '@store/store';

import React, { useState, useEffect } from 'react';

import { useKeyBindings } from 'rooks';

import { X } from '@ui/media/icons/X';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { UserX01 } from '@ui/media/icons/UserX01';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { useModKey } from '@shared/hooks/useModKey';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { TableInstance } from '@ui/presentation/Table';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { LinkedInSolid } from '@ui/media/icons/LinkedInSolid.tsx';
import { TableIdType, Organization, OrganizationStage } from '@graphql/types';
import { ActionItem } from '@organizations/components/Actions/ActionItem.tsx';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';
import { CreateContactFromLinkedInModal } from '@organizations/components/Actions/components/CreateContactFromLinkedInModal.tsx';

interface TableActionsProps {
  tableId?: TableIdType;
  focusedId?: string | null;
  onHide: (ids: string[]) => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<Store<Organization>>;
  onMerge: (primaryId: string, mergeIds: string[]) => void;
  onUpdateStage: (ids: string[], stage: OrganizationStage) => void;
  onCreateContact: (props: {
    socialUrl: string;
    organizationId: string;
    options?: {
      onSuccess?: (serverId: string) => void;
    };
  }) => void;
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
}: TableActionsProps) => {
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const {
    open: isCreateContactModalOpen,
    onOpen: onOpenCreateContactModal,
    onClose: onCloseCreateContactModal,
  } = useDisclosure({
    id: 'create-linkedin-contact',
  });
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

  const handleHideOrganizations = () => {
    onHide(selectedIds);
    onClose();
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
    if (!selectCount) return;
    onUpdateStage(selectedIds, OrganizationStage.Unqualified);
    clearSelection();
  };
  const moveToTarget = () => {
    if (!selectCount) return;
    onUpdateStage(selectedIds, OrganizationStage.Target);
    clearSelection();
  };
  const moveToOpportunities = () => {
    if (!selectCount) return;
    onUpdateStage(selectedIds, OrganizationStage.Engaged);
    clearSelection();
  };

  const createContactForOrganization = (url: string) => {
    if (!targetId) return;
    onCreateContact({
      socialUrl: url,
      organizationId: targetId || selectedIds[0],
      options: {
        onSuccess: () => {
          clearSelection();
        },
      },
    });
  };

  useKeyBindings(
    {
      u: moveToAllOrgs,
      t: moveToTarget,
      o: moveToOpportunities,
      l: () => tableId === TableIdType.Nurture && onOpenCreateContactModal(),
      Escape: clearSelection,
    },
    { when: enableKeyboardShortcuts },
  );

  useModKey(
    'v',
    () => {
      if (focusedId) {
        setTargetId(focusedId);
        onOpenCreateContactModal();
      }
    },
    {
      when:
        tableId === TableIdType.Nurture &&
        enableKeyboardShortcuts &&
        typeof focusedId === 'string',
    },
  );

  if (!selectCount && !targetId) {
    return null;
  }

  const getOrganizationName = () => {
    if (focusedId) {
      return table.getRow(focusedId)?.original?.value?.name || '';
    }
    if (targetId) {
      return table.getRow(targetId)?.original?.value?.name || '';
    }

    return '';
  };

  return (
    <>
      {selectCount > 0 && (
        <ButtonGroup className='flex items-center translate-x-[-50%] justify-center bottom-[32px] *:border-none'>
          {selectCount && (
            <div className='bg-gray-700 px-3 py-2 rounded-s-lg'>
              <p
                onClick={clearSelection}
                className='text-gray-25 text-sm font-semibold text-nowrap leading-5 outline-dashed outline-1 rounded-[2px] outline-gray-400 pl-2 pr-1 hover:bg-gray-800 transition-colors cursor-pointer'
              >
                {`${selectCount} selected`}
                <span className='ml-1'>
                  <X />
                </span>
              </p>
            </div>
          )}

          <ActionItem
            onClick={onOpen}
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
            [TableIdType.Leads, TableIdType.Nurture].includes(tableId) && (
              <ActionItem
                shortcutKey='U'
                onClick={moveToAllOrgs}
                tooltip={'Change to Unqualified and move to All orgs'}
                icon={<UserX01 className='text-inherit size-3' />}
              >
                Unqualify
              </ActionItem>
            )}

          {tableId && [TableIdType.Leads].includes(tableId) && (
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
            [TableIdType.Leads, TableIdType.Nurture].includes(tableId) && (
              <ActionItem
                shortcutKey='O'
                onClick={moveToOpportunities}
                tooltip='Change to Engaged and move to Opportunities'
                icon={<CoinsStacked01 className='text-inherit size-3' />}
              >
                Opportunity
              </ActionItem>
            )}
          {tableId &&
            selectedIds.length === 1 &&
            [TableIdType.Nurture].includes(tableId) && (
              <ActionItem
                shortcutKey='L'
                onClick={onOpenCreateContactModal}
                tooltip='Add contact via LinkedIn URL'
                icon={
                  <LinkedInSolid className='text-inherit size-4 text-inherit ' />
                }
              >
                Add LinkedIn contact
              </ActionItem>
            )}
        </ButtonGroup>
      )}

      <CreateContactFromLinkedInModal
        isOpen={isCreateContactModalOpen}
        onClose={onCloseCreateContactModal}
        organizationName={getOrganizationName()}
        onConfirm={createContactForOrganization}
      />
      <ConfirmDeleteDialog
        isOpen={isOpen}
        icon={<Archive />}
        onClose={onClose}
        confirmButtonLabel={'Archive'}
        onConfirm={handleHideOrganizations}
        loadingButtonLabel='Archiving'
        label={`Archive selected ${
          selectCount === 1 ? 'organization' : 'organizations'
        }?`}
      />
    </>
  );
};
