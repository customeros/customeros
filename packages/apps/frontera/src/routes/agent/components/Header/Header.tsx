import { Link } from 'react-router-dom';
import { useRef, useState, useEffect, useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Switch } from '@ui/form/Switch';
import { Icon, IconName } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

import { IconAndColorPicker } from './components/IconAndColorPicker';

interface HeaderProps {
  id: string;
  icon?: IconName;
  agentName: string;
  isActive: boolean;
  onToggleActive: () => void;
  colorMap: [ring: string, bg: string, iconColor: string];
}

export const Header = observer(
  ({
    onToggleActive,
    agentName,
    isActive,
    icon,
    id,
    colorMap,
  }: HeaderProps) => {
    const store = useStore();
    const [showIconPicker, setShowIconPicker] = useState(false);
    const [menuOpen, setMenuOpen] = useState(false);
    const iconPickerTriggerRef = useRef<HTMLButtonElement>(null);
    const shouldOpenIconPickerRef = useRef(false);

    useEffect(() => {
      if (shouldOpenIconPickerRef.current && !menuOpen) {
        shouldOpenIconPickerRef.current = false;

        const timer = setTimeout(() => {
          setShowIconPicker(true);
        }, 100);

        return () => clearTimeout(timer);
      }
    }, [menuOpen]);

    const handleIconPickerOpen = useCallback(() => {
      shouldOpenIconPickerRef.current = true;
      setMenuOpen(false);
    }, []);
    const [_ring, bg, iconColor] = colorMap;

    const handleOpenCommandMenu = (type: 'RenameAgent' | 'ArchiveAgent') => {
      store.ui.commandMenu.setType(type);
      store.ui.commandMenu.setContext({
        ...store.ui.commandMenu.context,
        ids: [id],
      });
      store.ui.commandMenu.setOpen(true);
    };

    return (
      <div className='w-full border-b border-b-gray-200 px-3 py-[11px] flex justify-between items-center'>
        <div className='flex items-center gap-1'>
          <Link
            to='/agents'
            className='text-md font-medium text-grayModern-500 hover:text-grayModern-700 transition-colors'
          >
            Agents
          </Link>
          <Icon name='chevron-right' className='w-4 h-4 text-gray-500' />
          <div className='flex items-center gap-1'>
            <Popover open={showIconPicker} onOpenChange={setShowIconPicker}>
              <PopoverTrigger className={'group'} ref={iconPickerTriggerRef}>
                <div
                  className={cn(
                    'flex items-center rounded-sm p-0.5 bg-grayModern-100',
                    bg,
                  )}
                >
                  <Icon
                    name={icon as IconName}
                    className={cn('size-4 text-grayModern-500', iconColor)}
                  />
                </div>
              </PopoverTrigger>

              <PopoverContent
                align='start'
                className='p-4 bg-grayModern-700 border-0'
              >
                <IconAndColorPicker />
              </PopoverContent>
            </Popover>

            <div className='flex items-center gap-1'>
              <Button
                size='xxs'
                variant='ghost'
                onClick={() => handleOpenCommandMenu('RenameAgent')}
                className='text-md font-medium hover:bg-transparent focus:bg-transparent'
              >
                {agentName}
              </Button>
              <Menu open={menuOpen} onOpenChange={setMenuOpen}>
                <MenuButton asChild>
                  <IconButton
                    size='xs'
                    variant='ghost'
                    aria-label={'Menu'}
                    icon={<Icon name='dots-vertical' />}
                  />
                </MenuButton>
                <MenuList side='bottom' align='start'>
                  <MenuItem
                    className='group'
                    onClick={() => handleOpenCommandMenu('RenameAgent')}
                  >
                    <Icon
                      name='edit-03'
                      className='text-grayModern-500 group-hover:text-grayModern-700'
                    />
                    Rename
                  </MenuItem>

                  <MenuItem className='group' onClick={handleIconPickerOpen}>
                    <Icon
                      name='brush-01'
                      className='text-grayModern-500 group-hover:text-grayModern-700'
                    />
                    Change icon
                  </MenuItem>
                  {/*<MenuItem*/}
                  {/*  disabled={true}*/}
                  {/*  className='group'*/}
                  {/*  onClick={() => {*/}
                  {/*    store.ui.commandMenu.setContext({*/}
                  {/*      ...store.ui.commandMenu.context,*/}
                  {/*      ids: [id],*/}
                  {/*    });*/}
                  {/*    store.ui.commandMenu.setType('DuplicateAgent');*/}
                  {/*    store.ui.commandMenu.setOpen(true);*/}
                  {/*  }}*/}
                  {/*>*/}
                  {/*  <Icon*/}
                  {/*    name='layers-two-01'*/}
                  {/*    className='text-grayModern-500 group-hover:text-grayModern-700'*/}
                  {/*  />*/}
                  {/*  Duplicate agent*/}
                  {/*</MenuItem>*/}
                  <MenuItem
                    className='group'
                    onClick={() => handleOpenCommandMenu('ArchiveAgent')}
                  >
                    <Icon
                      name='archive'
                      className='text-grayModern-500 group-hover:text-grayModern-700'
                    />
                    Archive agent
                  </MenuItem>
                </MenuList>
              </Menu>
            </div>

            <div className='ml-4 flex items-center'>
              <Switch checked={isActive} onChange={onToggleActive} />
            </div>
          </div>
        </div>
      </div>
    );
  },
);
