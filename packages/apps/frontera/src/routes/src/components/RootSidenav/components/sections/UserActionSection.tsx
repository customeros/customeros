import React from 'react';
import { useNavigate } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';
import { LogOut01 } from '@ui/media/icons/LogOut01';
import { Settings01 } from '@ui/media/icons/Settings01';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';

import { SidenavItem } from '../SidenavItem';

export const UserActionSection = () => {
  const navigate = useNavigate();
  const store = useStore();

  const handleSignOutClick = () => {
    store.session.clearSession();

    if (store.demoMode) {
      window.location.reload();

      return;
    }
    navigate('/auth/signin');
  };

  return (
    <div className='user-action-section'>
      <NotificationCenter />
      <SidenavItem
        label='Settings'
        isActive={false}
        dataTest='side-nav-item-settings'
        onClick={() => navigate('/settings')}
        icon={(isActive) => (
          <Settings01 className={isActive ? 'active-icon' : 'inactive-icon'} />
        )}
      />
      <SidenavItem
        label='Sign out'
        isActive={false}
        onClick={handleSignOutClick}
        dataTest='side-nav-item-sign-out'
        icon={(isActive) => (
          <LogOut01 className={isActive ? 'active-icon' : 'inactive-icon'} />
        )}
      />
    </div>
  );
};
