import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { Image } from '@ui/media/Image/Image';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import logoCustomerOs from '../../../../../src/assets/customer-os-small.png';

export const LogoSection = observer(() => {
  const store = useStore();
  const navigate = useNavigate();

  const handleSignOutClick = () => {
    store.session.clearSession();

    if (store.demoMode) {
      window.location.reload();

      return;
    }
    navigate('/auth/signin');
  };

  return (
    <div className='flex justify-between'>
      <Menu>
        <div className='py-2 pr-3 pl-[18px]'>
          <MenuButton
            data-test='logo-button'
            className='flex items-center gap-1.5 !outline-none'
          >
            <Image
              width={20}
              height={20}
              alt='CustomerOS'
              fallbackSrc={logoCustomerOs}
              className='logo-image rounded'
              src={store.settings.tenant.value?.workspaceLogo}
            />
            <span className='font-semibold  text-start w-[fit-content] overflow-hidden text-ellipsis whitespace-nowrap'>
              {store.settings.tenant.value?.workspaceName || 'CustomerOS'}
            </span>
            <Icon name='chevron-down' className='size-3 min-w-3' />
          </MenuButton>
        </div>
        <MenuList align='start' side='bottom' className='min-w-[137px]'>
          <MenuItem className='group' onClick={() => navigate('/settings')}>
            <div data-test='logo-settings' className='flex gap-2 items-center'>
              <Icon
                name='settings-02'
                className='group-hover:text-gray-700 text-gray-500'
              />
              <span>Settings</span>
            </div>
          </MenuItem>
          <MenuItem className='group' onClick={handleSignOutClick}>
            <div className='flex gap-2 items-center'>
              <Icon
                name='log-out-01'
                className='group-hover:text-gray-700 text-gray-500'
              />
              <span>Sign Out</span>
            </div>
          </MenuItem>
        </MenuList>
      </Menu>

      {(store.isSyncing || store.isBootstrapping) &&
        store.windowManager.networkStatus === 'online' && (
          <Tooltip
            side='bottom'
            delayDuration={0}
            label='Syncing the latest changes'
          >
            <div className='flex items-center'>
              <Spinner
                size='sm'
                label='Syncing'
                className='text-gray-300 fill-gray-700 mr-3'
              />
            </div>
          </Tooltip>
        )}

      {store.windowManager.networkStatus === 'offline' && (
        <Tooltip
          side='bottom'
          delayDuration={0}
          label={
            <>
              <p className='font-medium'>You’re offline</p>
              <p className='max-w-40'>
                Your changes are saving locally but may conflict with others
                when you reconnect
              </p>
            </>
          }
        >
          <div className='flex items-center'>
            <Icon name='cloud-off' className='mr-3 text-gray-500' />
          </div>
        </Tooltip>
      )}
    </div>
  );
});
