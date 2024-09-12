import React from 'react';
import { useNavigate } from 'react-router-dom';

import { Image } from '@ui/media/Image/Image';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { LogOut01 } from '@ui/media/icons/LogOut01';
import { Settings02 } from '@ui/media/icons/Settings02';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import logoCustomerOs from '../../../../../src/assets/customer-os-small.png';

export const LogoSection = () => {
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
    <Menu>
      <MenuButton className='py-2 px-3 pl-[18px] !outline-none'>
        <div data-test='logo-button' className='flex items-center gap-1.5'>
          {!isLoading ? (
            <>
              <Image
                width={20}
                height={20}
                alt='CustomerOS'
                className='logo-image rounded mt-0.5'
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
        </div>
      </MenuButton>
      <MenuList align='start' side='bottom' className='w-[180px] ml-2'>
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
  );
};
