import React from 'react';
import { SidePanelListItem } from './side-panel-list-item';
import styles from './side-panel.module.scss';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { useRecoilValue } from 'recoil';
import { logoutUrlState } from '../../../../state';
import { useJune } from '@spaces/hooks/useJune';
import { globalCacheData } from '../../../../state/globalCache';
import Portfolio from '@spaces/atoms/icons/Portfolio';
import Contacts from '@spaces/atoms/icons/Contacts';
import Company from '@spaces/atoms/icons/Company';
import Settings from '@spaces/atoms/icons/Settings';
import Exit from '@spaces/atoms/icons/Exit';
import Customer from '@spaces/atoms/icons/Customer';

export const SidePanel: React.FC = () => {
  const analytics = useJune();
  const router = useRouter();
  const logoutUrl = useRecoilValue(logoutUrlState);
  const { isOwner } = useRecoilValue(globalCacheData);

  return (
    <>
      <aside className={styles.sidebar}>
        <div className={styles.logoWrapper} role='button' tabIndex={0}>
          <Image
            src='/logos/openline_small.svg'
            alt='Openline'
            width={31}
            height={40}
            className={styles.logo}
          />
        </div>

        <SidePanelListItem
          label='Contacts'
          icon={<Contacts height={24} width={24} style={{ scale: '0.8' }} />}
          onClick={() => router.push('/contact')}
          selected={router.asPath.startsWith('/contact')}
        />
        <SidePanelListItem
          label='Organizations'
          icon={<Company height={24} width={24} style={{ scale: '0.8' }} />}
          onClick={() => router.push('/organization')}
          selected={
            router.asPath === '/' || router.asPath.startsWith('/organization')
          }
        />
        {isOwner && (
          <SidePanelListItem
            label='My portfolio'
            icon={<Portfolio height={24} width={24} style={{ scale: '0.8' }} />}
            onClick={() => router.push('/portfolio')}
            selected={
              router.asPath === '/' || router.asPath.startsWith('/portfolio')
            }
          />
        )}
        <SidePanelListItem
          label='Customers'
          icon={<Customer height={24} width={24} />}
          onClick={() => router.push('/customers')}
          selected={
            router.asPath === '/customers' ||
            router.asPath.startsWith('/customers')
          }
        />

        <div className={styles.bottom}>
          <SidePanelListItem
            label='Settings'
            icon={<Settings height={24} width={24} style={{ scale: '0.8' }} />}
            onClick={() => router.push('/settings')}
            selected={router.asPath.startsWith('/settings')}
          />
          <SidePanelListItem
            label='Log Out'
            icon={<Exit height={24} width={24} style={{ scale: '0.8' }} />}
            onClick={() => {
              document.cookie =
                'AUTH_CHECK=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
              analytics?.reset();
              window.location.href = logoutUrl;
            }}
          />
        </div>
      </aside>
    </>
  );
};
