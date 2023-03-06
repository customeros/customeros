import React, { useEffect, useRef } from 'react';

import { Building, Cog, SignOut, UserPlus } from '../../atoms';
import { SidePanelListItem } from './side-panel-list-item';
import classNames from 'classnames';
import styles from './side-panel.module.scss';
import Image from 'next/image';
import { useRouter } from 'next/router';
import { useRecoilValue } from 'recoil';
import { logoutUrlState } from '../../../../state';
import { useDetectClickOutside } from '../../../../hooks';

interface SidePanelProps {
  onPanelToggle: (status: boolean) => void;
  isPanelOpen: boolean;
  children: React.ReactNode;
}

export const SidePanel: React.FC<SidePanelProps> = ({
  onPanelToggle,
  isPanelOpen,
  children,
}) => {
  const router = useRouter();
  const logoutUrl = useRecoilValue(logoutUrlState);

  return (
    <>
      <aside
        className={classNames(styles.sidebar, {
          [styles.collapse]: !isPanelOpen,
        })}
      >
        <div
          className={styles.logoNameWrapper}
          role='button'
          tabIndex={0}
          onClick={() => onPanelToggle(!isPanelOpen)}
        >
          <Image
            src='logos/openline.svg'
            alt='Openline'
            width={120}
            height={40}
            className={styles.logoExpanded}
          />
          <Image
            src='logos/openline_small.svg'
            alt='Openline'
            width={40}
            height={40}
            className={styles.logoCollapsed}
          />
        </div>

        <ul className={styles.featuresList}>
          <SidePanelListItem
            label='Add organization'
            isOpen={isPanelOpen}
            icon={<Building />}
            onClick={() => router.push('/organization/new')}
          />
          <SidePanelListItem
            label='Add contact'
            isOpen={isPanelOpen}
            icon={<UserPlus />}
            onClick={() => router.push('/contact/new')}
          />
          <SidePanelListItem
            label='Settings'
            isOpen={isPanelOpen}
            icon={<Cog />}
            onClick={() => router.push('/settings')}
          />

          <SidePanelListItem
            label='Log Out'
            isOpen={isPanelOpen}
            icon={<SignOut />}
            onClick={() => (window.location.href = logoutUrl)}
          />
        </ul>
      </aside>
      {isPanelOpen && <div className={styles.overlay} onClick={() => onPanelToggle(!isPanelOpen)}/>}
      <div className={styles.webChat}>{children}</div>
    </>
  );
};
