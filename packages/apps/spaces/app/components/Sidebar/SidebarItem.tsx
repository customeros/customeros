'use client';
import { ReactNode } from 'react';
import Link from 'next/link';
import classNames from 'classnames';

import { Tooltip } from '@spaces/atoms/tooltip';

import styles from './SidebarItem.module.scss';
import { usePathname } from 'next/navigation';

interface SidebarItemProps {
  href?: string;
  label: string;
  icon?: ReactNode;
  onClick?: () => void;
}

export const SidebarItem = ({
  label,
  icon,
  href,
  onClick,
}: SidebarItemProps) => {
  const pathname = usePathname();
  const isActive = href ? pathname?.startsWith(href) : false;

  return (
    <div
      tabIndex={0}
      role='button'
      onClick={onClick}
      className={classNames(styles.featuresItem, {
        [styles.selected]: isActive,
      })}
    >
      <Link
        href={href ?? ''}
        className={styles.featuresItemIcon}
        id={`icon-${label}`}
      >
        {icon}
      </Link>
      <Tooltip
        content={label}
        showDelay={300}
        autoHide={false}
        position='right'
        target={`#icon-${label}`}
      />
    </div>
  );
};
