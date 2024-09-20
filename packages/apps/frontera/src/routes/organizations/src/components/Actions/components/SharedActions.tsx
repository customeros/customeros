import { ContactStore } from '@store/Contacts/Contact.store.ts';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';

import { X } from '@ui/media/icons/X.tsx';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Archive } from '@ui/media/icons/Archive.tsx';
import { TableInstance } from '@ui/presentation/Table';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';
import { ActionItem } from '@organizations/components/Actions/components/ActionItem.tsx';

interface TableActionsProps {
  onHide: () => void;
  selectCount: number;
  onOpenCommandK: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<ContactStore>;
  handleOpen: (type: CommandMenuType) => void;
}

export const SharedTableActions = ({
  table,
  onOpenCommandK,
  onHide,
  selectCount,
}: TableActionsProps) => {
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
          onClick={() => onHide()}
          dataTest='actions-archive'
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
        <ActionItem
          onClick={onOpenCommandK}
          dataTest='org-actions-commandk'
          icon={
            <span className='text-inherit w-auto h-auto'>
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
