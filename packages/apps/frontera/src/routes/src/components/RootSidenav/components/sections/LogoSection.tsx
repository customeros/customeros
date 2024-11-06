import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Image } from '@ui/media/Image/Image';
import { Spinner } from '@ui/feedback/Spinner';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { LogOut01 } from '@ui/media/icons/LogOut01';
import { CloudOff } from '@ui/media/icons/CloudOff';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Settings02 } from '@ui/media/icons/Settings02';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import logoCustomerOs from '../../../../../src/assets/customer-os-small.png';

export const LogoSection = observer(() => {
  const store = useStore();
  const navigate = useNavigate();

  const isLoading = store.globalCache?.isLoading;

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
        <div data-test='logo-button' className='py-2 pr-3 pl-[18px]'>
          <MenuButton className='flex items-center gap-1.5 !outline-none'>
            {!isLoading ? (
              <>
                <Image
                  width={20}
                  height={20}
                  alt='CustomerOS'
                  className='logo-image rounded'
                  src={
                    store.settings.tenant.value?.workspaceLogo || logoCustomerOs
                  }
                />
                <span className='font-semibold  text-start w-[fit-content] overflow-hidden text-ellipsis whitespace-nowrap'>
                  {store.settings.tenant.value?.workspaceName || 'CustomerOS'}
                </span>
                <ChevronDown className='size-3 min-w-3' />
              </>
            ) : (
              <Skeleton className='w-full h-8 mr-2' />
            )}
          </MenuButton>
        </div>
        <MenuList align='start' side='bottom' className='min-w-[137px]'>
          <MenuItem className='group' onClick={() => navigate('/settings')}>
            <div className='flex gap-2 items-center'>
              <Settings02 className='group-hover:text-gray-700 text-gray-500' />
              <span>Settings</span>
            </div>
          </MenuItem>
          <MenuItem className='group' onClick={handleSignOutClick}>
            <div className='flex gap-2 items-center'>
              <LogOut01 className='group-hover:text-gray-700 text-gray-500' />
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
            <CloudOff className='mr-3 text-gray-500' />
          </div>
        </Tooltip>
      )}
    </div>
  );
});
