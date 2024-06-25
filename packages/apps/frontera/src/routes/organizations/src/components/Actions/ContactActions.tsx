import type { Store } from '@store/store';

import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { X } from '@ui/media/icons/X';
import { Contact } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { TableInstance } from '@ui/presentation/Table';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ActionItem } from '@organizations/components/Actions/ActionItem.tsx';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

interface TableActionsProps {
  table: TableInstance<Store<Contact>>;
}

export const ContactTableActions = observer(({ table }: TableActionsProps) => {
  const store = useStore();

  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const [targetId, setTargetId] = useState<string | null>(null);

  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection);
  const selectCount = selectedIds.length;

  const clearSelection = () => table.resetRowSelection();

  const handleHideContacts = () => {
    selectedIds.forEach((id) => {
      store.contacts.remove(id);
    });
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
      </ButtonGroup>

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
});
