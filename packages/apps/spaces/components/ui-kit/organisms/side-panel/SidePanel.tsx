import React from 'react';
import { SidePanelListItem } from './side-panel-list-item';
import styles from './side-panel.module.scss';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { useRecoilValue } from 'recoil';
import { useJune } from '@spaces/hooks/useJune';
import { globalCacheData } from '../../../../state/globalCache';
import Portfolio from '@spaces/atoms/icons/Portfolio';
import Company from '@spaces/atoms/icons/Company';
import Settings from '@spaces/atoms/icons/Settings';
import Exit from '@spaces/atoms/icons/Exit';
import Customer from '@spaces/atoms/icons/Customer';
import { signOut } from 'next-auth/react';

export const SidePanel: React.FC = () => {
  const analytics = useJune();
  const router = useRouter();
  // const { isOwner } = useRecoilValue(globalCacheData);

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
          label='Organizations'
          icon={<Company height={24} width={24} style={{ scale: '0.8' }} />}
          onClick={() => router.push('/organization')}
          selected={
            router.asPath === '/' || router.asPath.startsWith('/organization')
          }
        />
        <SidePanelListItem
          label='Customers'
          icon={<Customer height={24} width={24} style={{ scale: '0.8' }} />}
          onClick={() => router.push('/customers')}
          selected={
            router.asPath === '/customers' ||
            router.asPath.startsWith('/customers')
          }
        />
        {/*{isOwner && (*/}
        {/*  <SidePanelListItem*/}
        {/*    label='My portfolio'*/}
        {/*    icon={<Portfolio height={24} width={24} style={{ scale: '0.8' }} />}*/}
        {/*    onClick={() => router.push('/portfolio')}*/}
        {/*    selected={*/}
        {/*      router.asPath === '/portfolio' ||*/}
        {/*      router.asPath.startsWith('/portfolio')*/}
        {/*    }*/}
        {/*  />*/}
        {/*)}*/}

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
              analytics?.reset();
              signOut();
            }}
          />
        </div>
      </aside>
    </>
  );
};
