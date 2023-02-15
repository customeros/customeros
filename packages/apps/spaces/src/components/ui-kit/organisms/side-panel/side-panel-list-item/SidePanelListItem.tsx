import React, { EventHandler, ReactNode } from 'react';
import styles from './side-panel-list-item.module.scss';
import classNames from 'classnames';
import { Button } from '../../../atoms';
interface SidePanelListItemProps {
  label: string;
  icon?: ReactNode;
  isOpen: boolean;
  onClick: EventHandler<never>;
}

export const SidePanelListItem: React.FC<SidePanelListItemProps> = ({
  label,
  icon,
  isOpen,
  onClick,
}) => {
  return (
    <li className={styles.featuresItem} role='button' tabIndex={0}>
      {icon && <span className={styles.featuresItemIcon}>{icon}</span>}

      <span
        className={classNames(styles.featuresItemText, {
          [styles.featuresItemTextHidden]: !isOpen,
        })}
      >
        {label}
      </span>
    </li>
  );
};
