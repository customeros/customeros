import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Icon } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Devtools } from '@shared/components/Devtools/Devtools';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

export const SettingsSection = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const debuggerFlag = useFeatureIsOn('debugger');
  const { open, onOpen, onClose, onToggle } = useDisclosure();
  const isDebuggerEnabled = import.meta.env.DEV || debuggerFlag;

  const handleSignOutClick = () => {
    store.session.clearSession();

    if (store.demoMode) {
      window.location.reload();

      return;
    }
    navigate('/auth/signin');
  };

  const navigateToDocs = (url: string) => {
    window.open(url, '_blank', 'noopener,noreferrer');
  };

  const handleOpenAtlas = () => {
    document.getElementById('atlas-launcher')?.click?.();
  };

  return (
    <div className='flex justify-between items-center px-2'>
      <div className='flex items-center gap-x-2'>
        <Tooltip label='Settings'>
          <IconButton
            size={'xs'}
            variant='ghost'
            aria-label={'Settings'}
            dataTest={'settings-button'}
            icon={<Icon name={'settings-02'} />}
            onClick={() => navigate('/settings')}
          />
        </Tooltip>

        <Menu>
          <Tooltip label='Help me'>
            <MenuButton
              asChild
              data-test='help-button'
              className='flex items-center gap-1.5 !outline-none'
            >
              <IconButton
                size={'xs'}
                variant='ghost'
                aria-label={'Help'}
                icon={<Icon name={'help-circle'} />}
              />
            </MenuButton>
          </Tooltip>

          <MenuList side='top' align='start' className='w-[184px]'>
            <MenuItem
              className='group'
              data-test='help-item-settings'
              onClick={() => store.ui.setShortcutsPanel(true)}
            >
              <div
                data-test='logo-settings'
                className='flex gap-2 items-center'
              >
                <Icon
                  name='keyboard-02'
                  className='group-hover:text-gray-700 text-gray-500'
                />
                <span>Shortcuts</span>
              </div>
            </MenuItem>
            <MenuItem
              className='group'
              data-test='help-item-nav-api-docs'
              onClick={() =>
                navigateToDocs('https://docs.customeros.ai/api-overview')
              }
            >
              <div className='flex w-full justify-between items-center'>
                <div className='flex gap-2 items-center'>
                  <Icon
                    name='code-browser'
                    className='group-hover:text-gray-700 text-gray-500'
                  />
                  <span>API docs</span>
                </div>
                <Icon
                  name='arrow-narrow-up-right'
                  className='group-hover:text-gray-500 text-gray-400'
                />
              </div>
            </MenuItem>
            <MenuItem
              className='group'
              data-test='help-item-nav-docs'
              onClick={() =>
                navigateToDocs('https://docs.customeros.ai/introduction')
              }
            >
              <div className='flex w-full justify-between items-center'>
                <div className='flex gap-2 items-center'>
                  <Icon
                    name='book-closed'
                    className='group-hover:text-gray-700 text-gray-500'
                  />
                  <span>Help docs</span>
                </div>
                <Icon
                  name='arrow-narrow-up-right'
                  className='group-hover:text-gray-500 text-gray-400'
                />
              </div>
            </MenuItem>
            <MenuItem
              className='group'
              onClick={handleOpenAtlas}
              data-test='ask-for-help-button'
            >
              <div className='flex gap-2 items-center'>
                <Icon
                  name='message-smile-square'
                  className='group-hover:text-gray-700 text-gray-500'
                />
                <span>Ask for help</span>
              </div>
            </MenuItem>

            {isDebuggerEnabled && (
              <MenuItem onClick={onOpen} className='group'>
                <Icon
                  name='code-square-02'
                  className='group-hover:text-gray-700 text-gray-500'
                />
                <span>Debugger</span>
              </MenuItem>
            )}
          </MenuList>
        </Menu>
      </div>
      <Tooltip label={'Sign out'}>
        <IconButton
          size={'xs'}
          variant='ghost'
          aria-label={'Log out'}
          dataTest={'sign-out-button'}
          onClick={handleSignOutClick}
          icon={<Icon name={'log-out-01'} />}
        />
      </Tooltip>

      {isDebuggerEnabled && (
        <Devtools open={open} onClose={onClose} onToggle={onToggle} />
      )}
    </div>
  );
});
