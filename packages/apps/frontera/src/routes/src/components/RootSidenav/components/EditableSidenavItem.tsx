import React, { useState, ReactElement, MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive.tsx';
import { TextInput } from '@ui/media/icons/TextInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
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
        ? ['font-semibold', 'bg-gray-100']
        : ['font-normal', 'bg-transparent'],
    );

    return (
      <Button
        size='md'
        variant='ghost'
        colorScheme='gray'
        data-test={dataTest}
        onClick={handleClick}
        className={`w-full justify-start px-3 text-gray-700 group focus:shadow-EditableSideNavItemFocus ${dynamicClasses}`}
      >
        <div>{typeof icon === 'function' ? icon(!!isActive) : icon}</div>
        <div
          className={cn(
            'w-full flex overflow-ellipsis group-hover:overflow-hidden group-focus:overflow-hidden',
            {
              'overflow-hidden': isEditing,
            },
          )}
        >
          {label}
        </div>

        <div
          className={cn(
            'justify-end opacity-0 group-hover:opacity-100 group-focus:opacity-100 ',
            {
              'opacity-100': isEditing,
            },
          )}
        >
          <Menu open={isEditing} onOpenChange={setIsEditing}>
            <Tooltip label={label}>
              <MenuButton className='min-w-6 h-5 rounded-md outline-none focus:outline-none text-gray-400 hover:text-gray-500'>
                <DotsVertical className='text-inherit' />
              </MenuButton>
            </Tooltip>

            <MenuList align='end' side='bottom'>
              <MenuItem
                className='py-2.5'
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  store.ui.commandMenu.setType('RenameTableViewDef');
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                <TextInput className='text-gray-500' />
                Rename view
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
                <Archive />
                Archive view
              </MenuItem>
            </MenuList>
          </Menu>
        </div>
      </Button>
    );
  },
);
