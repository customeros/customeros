import React from 'react';

import { Cog, SignOut } from '../../atoms';
import { SidePanelListItem } from './side-panel-list-item';
import classNames from 'classnames';
import styles from './side-panel.module.scss';
import Image from 'next/image';

interface SidePanelProps {
  onPanelToggle: (status: boolean) => void;
  isPanelOpen: boolean;
  logoutUrl: string | undefined;
  children: React.ReactNode;
}

export const SidePanel: React.FC<SidePanelProps> = ({
  onPanelToggle,
  isPanelOpen,
  children,
  logoutUrl,
}) => {
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
            label='Settings'
            isOpen={isPanelOpen}
            icon={<Cog />}
            onClick={() => null}
          />
          <SidePanelListItem
            label='Log Out'
            isOpen={isPanelOpen}
            icon={<SignOut />}
            onClick={() => (window.location.href = logoutUrl ?? '#')}
          />
        </ul>
      </aside>
      <div
        className={styles.webChat}
        // @ts-expect-error this is valid syntax
        style={{ '--web-chat-left': isPanelOpen ? '135px' : '20px' }}
      >
        {children}
      </div>
    </>
  );
};
