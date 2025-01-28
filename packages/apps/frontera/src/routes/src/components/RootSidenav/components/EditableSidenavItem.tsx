import { useState, type ReactElement, type MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { buttonSize } from '@ui/form/Button/Button';
import { TextInput } from '@ui/media/icons/TextInput';
import { ghostButton } from '@ui/form/Button/Button.variants';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

interface EditableSideNavItemProps {
  id: string;
  href?: string;
  label: string;
  dataTest?: string;
  isActive?: boolean;
  onClick?: () => void;
  icon: ((isActive: boolean) => ReactElement) | ReactElement;
}

export const EditableSideNavItem = observer(
  ({
    label,
    icon,
    onClick,
    isActive,
    dataTest,
    id,
  }: EditableSideNavItemProps) => {
    const store = useStore();
    const [isEditing, setIsEditing] = useState(false);

    const handleClick: MouseEventHandler = (e) => {
      e.preventDefault();
      onClick?.();
    };

    const dynamicClasses = cn(
      isActive
        ? ['font-medium', 'bg-grayModern-100']
        : ['font-normal', 'bg-transparent'],
    );

    return (
      <div
        data-test={dataTest}
        onClick={handleClick}
        className={cn(
          buttonSize({ size: 'sm' }),
          (ghostButton({ colorScheme: 'grayModern' }),
          `flex w-full justify-start items-center gap-2 px-3 text-gray-700 cursor-pointer hover:bg-grayModern-100 *:hover:text-gray-700  group focus:shadow-EditableSideNavItemFocus mb-[2px] rounded-md ${dynamicClasses}`),
        )}
      >
        <div className='mt-[-1px]'>
          {typeof icon === 'function' ? icon(isActive!) : icon}
        </div>
        <div
          className={cn(
            'w-full text-justify whitespace-nowrap overflow-hidden overflow-ellipsis',
          )}
        >
          {label}
        </div>

        <div
          className={cn(
            'justify-end opacity-0 w-0 group-hover:opacity-100 group-focus:opacity-100 group-hover:w-6 group-focus:w-6',
            {
              'opacity-100 w-6': isEditing,
            },
          )}
        >
          <Menu open={isEditing} onOpenChange={setIsEditing}>
            <MenuButton className='min-w-6 h-5 rounded-md outline-none focus:outline-none text-gray-400 hover:text-gray-500 flex items-center'>
              <Icon name='dots-vertical' className='text-inherit' />
            </MenuButton>

            <MenuList align='end' side='bottom'>
              <MenuItem
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();

                  store.ui.commandMenu.setContext({
                    ids: [id],
                    entity: 'TableViewDef',
                  });
                  store.ui.commandMenu.setType('RenameTableViewDef');
                  store.ui.commandMenu.setOpen(true);
                  setTimeout(() => {
                    setIsEditing(false);
                  }, 1);
                }}
              >
                <TextInput className='text-gray-500' />
                Rename
              </MenuItem>
              <MenuItem
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  store.ui.commandMenu.toggle('DuplicateView');
                  store.ui.commandMenu.setContext({
                    ids: [id],
                    entity: 'TableViewDef',
                  });
                  // store.tableViewDefs.createFavorite(preset);
                  setIsEditing(false);
                }}
              >
                <Icon name='layers-two-01' className='text-gray-500' />
                Save as...
              </MenuItem>
              <MenuItem
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();

                  store.ui.commandMenu.setContext({
                    ids: [id],
                    entity: 'TableViewDef',
                  });
                  store.ui.commandMenu.setType('DeleteConfirmationModal');
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                <Icon name='archive' className='text-gray-500' />
                Archive view
              </MenuItem>
            </MenuList>
          </Menu>
        </div>
      </div>
    );
  },
);
