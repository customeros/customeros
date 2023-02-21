import React, { EventHandler, ReactNode } from 'react';
import styles from './side-panel-list-item.module.scss';
import classNames from 'classnames';
interface SidePanelListItemProps {
  label: string;
  icon?: ReactNode;
  isOpen: boolean;
  onClick: EventHandler<never>;
}
classNames(styles.featuresItemText, {});
export const SidePanelListItem: React.FC<SidePanelListItemProps> = ({
  label,
  icon,
  isOpen,
  onClick,
}) => {
  return (
    <li
      className={classNames(styles.featuresItem, {
        [styles.featuresItemHidden]: !isOpen,
      })}
      role='button'
      tabIndex={0}
      onClick={onClick}
    >
      {icon && <span className={styles.featuresItemIcon}>{icon}</span>}
      <span className={styles.featuresItemText}>{label}</span>
    </li>
  );
};
