import Image from 'next/image';

import Company from '@spaces/atoms/icons/Company';
import Settings from '@spaces/atoms/icons/Settings';
import Contacts from '@spaces/atoms/icons/Contacts';
import Customer from '@spaces/atoms/icons/Customer';
import Portfolio from '@spaces/atoms/icons/Portfolio';

import { SidebarItem } from './SidebarItem';
import { LogoutSidebarItem } from './LogoutSidebarItem';

import styles from './Sidebar.module.scss';

interface SidebarProps {
  isOwner: boolean;
}

export const Sidebar = ({ isOwner }: SidebarProps) => {
  return (
    <aside className={styles.sidebar}>
      <div className={styles.logoWrapper} role='button' tabIndex={0}>
        <Image
          width={31}
          height={40}
          alt='Openline'
          className={styles.logo}
          src='/logos/openline_small.svg'
        />
      </div>

      <SidebarItem
        href='/organization'
        label='Organizations'
        icon={<Company height={24} width={24} style={{ scale: '0.8' }} />}
      />
      <SidebarItem
        href='/customers'
        label='Customers'
        icon={<Customer height={24} width={24} style={{ scale: '0.8' }} />}
      />
      {isOwner && (
        <SidebarItem
          href='/portfolio'
          label='My portfolio'
          icon={<Portfolio height={24} width={24} style={{ scale: '0.8' }} />}
        />
      )}

      <SidebarItem
        href='/contact'
        label='Contacts'
        icon={<Contacts height={24} width={24} style={{ scale: '0.8' }} />}
      />

      <div className={styles.bottom}>
        <SidebarItem
          href='/settings'
          label='Settings'
          icon={<Settings height={24} width={24} style={{ scale: '0.8' }} />}
        />
        <LogoutSidebarItem />
      </div>
    </aside>
  );
};
