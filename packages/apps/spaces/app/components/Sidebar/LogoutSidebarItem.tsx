'use client';
import Exit from '@spaces/atoms/icons/Exit';
import { useJune } from '@spaces/hooks/useJune';

import { SidebarItem } from './SidebarItem';

export const LogoutSidebarItem = () => {
  const analytics = useJune();
  const logoutUrl =
    typeof window !== 'undefined'
      ? window?.sessionStorage?.getItem('logout_url')
      : null;

  const handleClick = () => {
    analytics?.reset();
    if (logoutUrl) {
      window.location.href = logoutUrl;
    }
  };

  return (
    <SidebarItem
      label='Logout'
      onClick={handleClick}
      icon={<Exit height={24} width={24} style={{ scale: '0.8' }} />}
    />
  );
};
