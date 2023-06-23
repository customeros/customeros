import React from 'react';
import classNames from 'classnames';
import styles from '@spaces/ui/form/select/components/select.module.scss';
import indicatorStyles from './health-selector.module.scss';
import { useSingleSelect } from '@spaces/ui/form/select/components/single-select/SingleSelect';

export const HealthIndicatorSelectMenu = () => {
  const { state, getMenuProps, getMenuItemProps } = useSingleSelect();

  const maxMenuHeight = 28 * 6;


  return (
    <ul
      className={styles.dropdownMenu}
      {...getMenuProps({ maxHeight: maxMenuHeight })}
    >
      {state.items.map(({ value, label }, index) => (
        <li
          key={value}
          className={classNames(styles.dropdownMenuItem, {
            [styles.isFocused]: state.currentIndex === index,
            [styles.isSelected]: state.selection === value,
          })}
          {...getMenuItemProps({ value, index })}
        >
          <div
            className={indicatorStyles.colorIndicator}
            style={{
              background: `var(--health-indicator-${label.toLowerCase()})`,
            }}
          />
          {label}
        </li>
      ))}
    </ul>
  );
};
