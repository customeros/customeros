import React, { EventHandler, ReactNode } from 'react';
import styles from './side-panel-list-item.module.scss';
import classNames from 'classnames';
import { Tooltip } from '@spaces/atoms/tooltip';
import { uuidv4 } from '../../../../../utils';

interface SidePanelListItemProps {
  label: string;
  icon?: ReactNode;
  onClick: EventHandler<never>;
  selected?: boolean;
}
classNames(styles.featuresItemText, {});
export const SidePanelListItem: React.FC<SidePanelListItemProps> = ({
  label,
  icon,
  onClick,
  selected,
}) => {
  const id = uuidv4();
  return (
    <div
      className={classNames(styles.featuresItem, {
        [styles.selected]: selected,
      })}
      role='button'
      tabIndex={0}
      onClick={onClick}
    >
      <div className={styles.featuresItemIcon} id={`side-panel-item-${id}`}>
        {icon}
      </div>
      <Tooltip
        content={label}
        target={`#side-panel-item-${id}`}
        position='right'
        showDelay={300}
        autoHide={false}
      />
    </div>
  );
};
