import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';

import { X } from '@ui/media/icons/X.tsx';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Archive } from '@ui/media/icons/Archive.tsx';
import { TableInstance } from '@ui/presentation/Table';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';
import { ActionItem } from '@organizations/components/Actions/components/ActionItem.tsx';

interface TableActionsProps {
  onHide: () => void;
  onOpenCommandK: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<ContactStore>;
  handleOpen: (type: CommandMenuType) => void;
}

export const SharedTableActions = ({
  table,
  onOpenCommandK,
  onHide,
}: TableActionsProps) => {
  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection);
  const selectCount = selectedIds.length;
  const clearSelection = () => table.resetRowSelection();

  if (!selectCount) return null;

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
          tooltip='Archive'
          onClick={() => onHide()}
          dataTest='contacts-actions-archive'
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
