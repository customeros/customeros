'use client';
import { signOut } from 'next-auth/react';

import Exit from '@spaces/atoms/icons/Exit';
import { useJune } from '@spaces/hooks/useJune';

import { SidebarItem } from './SidebarItem';

export const LogoutSidebarItem = () => {
  const analytics = useJune();

  const handleClick = () => {
    analytics?.reset();
    signOut();
  };

  return (
    <SidebarItem
      label='Logout'
      onClick={handleClick}
      icon={<Exit height={24} width={24} style={{ scale: '0.8' }} />}
    />
  );
};
