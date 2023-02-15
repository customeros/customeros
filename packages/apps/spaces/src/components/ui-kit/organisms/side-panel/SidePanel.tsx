import React from 'react';
import { IconButton } from '../../atoms/icon-button';
import {
  Button,
  ChevronLeft,
  ChevronRight,
  Cog,
  OpenlineLogo,
  SignOut,
} from '../../atoms';
import { SidePanelListItem } from './side-panel-list-item';
import classNames from 'classnames';
import styles from './side-panel.module.scss';
import Image from 'next/image';

interface SidePanelProps {
  onPanelToggle: (status: boolean) => void;
  isPanelOpen: boolean;
}

export const SidePanel: React.FC<SidePanelProps> = ({
  onPanelToggle,
  isPanelOpen,
}) => {
  return (
    <aside
      className={classNames(styles.sidebar, {
        [styles.collapse]: !isPanelOpen,
      })}
    >
      <div className={styles.logoNameWrapper}>
        <Image
          src='logos/openline.svg'
          alt='Openline'
          width={120}
          height={40}
          className={styles.logoExpanded}
        />
        <Image
          src='icons/openlineLogo.svg'
          alt='Openline'
          width={40}
          height={40}
          className={styles.logoCollapsed}
        />
      </div>

      <IconButton
        mode='secondary'
        className={styles.collapseExpandButton}
        onClick={() => onPanelToggle(!isPanelOpen)}
        icon={isPanelOpen ? <ChevronLeft /> : <ChevronRight />}
      />

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
          onClick={() => null}
        />
      </ul>
    </aside>
  );
};
