import { useState, useEffect } from 'react';

import { useKeys, useKeyBindings } from 'rooks';
import { ContactStore } from '@store/Contacts/Contact.store';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';

import { Tag } from '@graphql/types';
import { X } from '@ui/media/icons/X';
import { Archive } from '@ui/media/icons/Archive';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { TableInstance } from '@ui/presentation/Table';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';
import { ActionItem } from '@organizations/components/Actions/ActionItem.tsx';

interface TableActionsProps {
  focusedId: string | null;
  onOpenCommandK: () => void;
  onHideContacts: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<ContactStore>;
  handleOpen: (type: CommandMenuType) => void;
  onAddTags: (ids: string[], tags: Tag[]) => void;
}

export const ContactTableActions = ({
  table,
  enableKeyboardShortcuts,
  onOpenCommandK,
  handleOpen,
  onHideContacts,
  focusedId,
}: TableActionsProps) => {
  const [targetId, setTargetId] = useState<string | null>(null);

  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection);
  const selectCount = selectedIds.length;

  const clearSelection = () => table.resetRowSelection();

  useEffect(() => {
    if (selectCount === 1) {
      setTargetId(selectedIds[0]);
    }

    if (selectCount < 1) {
      setTargetId(null);
    }
  }, [selectCount]);

  useKeys(
    ['Shift', 'T'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('EditPersonaTag');
    },
    { when: enableKeyboardShortcuts },
  );
  useKeys(
    ['Shift', 'E'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('EditEmail');
    },
    { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
  );

  useKeys(
    ['Shift', 'R'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('EditName');
    },
    { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
  );
  useKeyBindings(
    {
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
          onClick={() => onHideContacts()}
          icon={<Archive className='text-inherit size-3' />}
        >
          Archive
        </ActionItem>
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
    </>
  );
};
