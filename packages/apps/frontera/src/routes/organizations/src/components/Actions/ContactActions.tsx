import type { Store } from '@store/store';

import { useState, useEffect } from 'react';

import { useKeyBindings } from 'rooks';

import { X } from '@ui/media/icons/X';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { Archive } from '@ui/media/icons/Archive';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { User02 } from '@ui/media/icons/User02.tsx';
import { Tags } from '@organization/components/Tabs';
import { TableInstance } from '@ui/presentation/Table';
import { Tag, Contact, DataSource } from '@graphql/types';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ActionItem } from '@organizations/components/Actions/ActionItem.tsx';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

interface TableActionsProps {
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<Store<Contact>>;
  onHideContacts: (ids: string[]) => void;
  onAddTags: (ids: string[], tags: Tag[]) => void;
}

export const ContactTableActions = ({
  table,
  enableKeyboardShortcuts,
  onAddTags,
  onHideContacts,
}: TableActionsProps) => {
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const {
    open: isTagEditOpen,
    onOpen: onOpenTagEdit,
    onClose: onCloseTagEdit,
  } = useDisclosure({ id: 'contact-actions' });
  const [targetId, setTargetId] = useState<string | null>(null);
  const [selectedTags, setSelectedTags] = useState<
    Array<{ label: string; value: string }>
  >([]);

  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection);
  const selectCount = selectedIds.length;

  const clearSelection = () => table.resetRowSelection();

  const handleHideContacts = () => {
    onHideContacts(selectedIds);
    setSelectedTags([]);
    onClose();
    clearSelection();
  };

  const handleBulkTagEditModal = () => {
    setSelectedTags([]);
    onCloseTagEdit();
  };

  const handleAddTags = () => {
    const tags = selectedTags.map((e) => ({
      name: e.label,
      id: e.value,
      appSource: 'organization',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      source: DataSource.Openline,
      metadata: {
        id: e.value,
        source: DataSource.Openline,
        sourceOfTruth: DataSource.Openline,
        appSource: 'organization',
        created: new Date().toISOString(),
        lastUpdated: new Date().toISOString(),
      },
    }));
    onAddTags(selectedIds, tags);
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

  useKeyBindings(
    {
      p: onOpenTagEdit,
      Escape: clearSelection,
    },
    { when: enableKeyboardShortcuts },
  );

  if (!selectCount && !targetId) return null;

  return (
    <>
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
        <ActionItem
          onClick={onOpenTagEdit}
          tooltip='Edit persona tags'
          shortcutKey='P'
          icon={<User02 className='text-inherit size-3' />}
        >
          Edit persona
        </ActionItem>
      </ButtonGroup>

      <ConfirmDeleteDialog
        isOpen={isTagEditOpen}
        icon={<Tag01 />}
        onClose={handleBulkTagEditModal}
        confirmButtonLabel={'Add tags'}
        onConfirm={handleAddTags}
        colorScheme='primary'
        loadingButtonLabel='Adding tags'
        label={`Add persona tags to ${selectCount} ${
          selectCount === 1 ? 'contact' : 'contacts'
        }?`}
        body={
          <div>
            <p className='text-gray-700 text-base font-normal mb-5'>
              What persona tags would you like to add to your selected contacts?
            </p>
            <Tags
              autofocus
              placeholder='Persona'
              icon={null}
              closeMenuOnSelect={true}
              value={selectedTags}
              onChange={(e) => setSelectedTags(e)}
            />
          </div>
        }
      />

      <ConfirmDeleteDialog
        isOpen={isOpen}
        icon={<Archive />}
        onClose={onClose}
        confirmButtonLabel={'Archive'}
        onConfirm={handleHideContacts}
        loadingButtonLabel='Archiving'
        label={`Archive selected ${
          selectCount === 1 ? 'contact' : 'contacts'
        }?`}
      />
    </>
  );
};
