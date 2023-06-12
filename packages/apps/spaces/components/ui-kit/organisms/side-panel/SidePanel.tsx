import React from 'react';

import { Building, Cog, SignOut, UserPlus } from '../../atoms/icons';
import { SidePanelListItem } from './side-panel-list-item';
import classNames from 'classnames';
import styles from './side-panel.module.scss';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { useRecoilValue } from 'recoil';
import { logoutUrlState } from '../../../../state';
import { useJune } from '@spaces/hooks/useJune';
import User from '@spaces/atoms/icons/User';

export const SidePanel: React.FC = () => {
  const analytics = useJune();
  const router = useRouter();
  const logoutUrl = useRecoilValue(logoutUrlState);
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
          icon={<User height={24} width={24} />}
          onClick={() => router.push('/contact')}
          selected={router.asPath.startsWith('/contact')}
        />
        <SidePanelListItem
          label='Organizations'
          icon={<Building height={24} width={24} />}
          onClick={() => router.push('/organization')}
          selected={
            router.asPath === '/' || router.asPath.startsWith('/organization')
          }
        />
        {/*<SidePanelListItem*/}
        {/*  label='My portfolio'*/}
        {/*  icon={<Building height={24} width={24} />}*/}
        {/*  onClick={() => router.push('/portfolio')}*/}
        {/*  selected={*/}
        {/*    router.asPath === '/' || router.asPath.startsWith('/portfolio')*/}
        {/*  }*/}
        {/*/>*/}

        <div className={styles.bottom}>
          <SidePanelListItem
            label='Settings'
            icon={<Cog height={24} width={24} />}
            onClick={() => router.push('/settings')}
            selected={router.asPath.startsWith('/settings')}
          />
          <SidePanelListItem
            label='Log Out'
            icon={<SignOut height={24} width={24} />}
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
