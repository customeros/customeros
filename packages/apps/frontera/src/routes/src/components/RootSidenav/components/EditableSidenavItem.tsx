import { useSearchParams } from 'react-router-dom';
import React, { useState, ReactElement, MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive.tsx';
import { TextInput } from '@ui/media/icons/TextInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { LayersTwo01 } from '@ui/media/icons/LayersTwo01.tsx';
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
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset') ?? '1';
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
            'w-full text-justify overflow-hidden overflow-ellipsis',
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
            <Tooltip label={label}>
              <MenuButton className='min-w-6 h-5 rounded-md outline-none focus:outline-none text-gray-400 hover:text-gray-500'>
                <DotsVertical className='text-inherit' />
              </MenuButton>
            </Tooltip>

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
                  setIsEditing(false);
                }}
              >
                <TextInput className='text-gray-500' />
                Rename view
              </MenuItem>
              <MenuItem
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  store.tableViewDefs.createFavorite(preset);
                  setIsEditing(false);
                }}
              >
                <LayersTwo01 className='text-gray-500' />
                Duplicate view
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
                <Archive className='text-gray-500' />
                Archive view
              </MenuItem>
            </MenuList>
          </Menu>
        </div>
      </Button>
    );
  },
);
