import classNames from 'classnames';
import styles from './select.module.scss';
import React from 'react';
import { useSelect } from '../index';

interface SelectMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
  isCreatable?: boolean;
}

export const CreatableListItem = () => {
  const { state, getMenuItemProps } = useSelect();

  return (
    <li
      key='new-item'
      className={classNames(styles.dropdownMenuItem, {
        [styles.isFocused]: state.currentIndex === 0,
        [styles.isSelected]: true,
      })}
      {...getMenuItemProps({ value: state.value, index: 0 })}
    >
      Create new option &quot;{state.value}&quot;
    </li>
  );
};

export const CreatableSelectMenu = ({
  noOfVisibleItems = 7,
  itemSize = 38,
}: SelectMenuProps) => {
  const { state, getMenuProps, getMenuItemProps } = useSelect();
  const maxMenuHeight = itemSize * noOfVisibleItems;
  const showCreateOption = state.value.length;

  return (
    <ul
      className={styles.dropdownMenu}
      {...getMenuProps({ maxHeight: maxMenuHeight })}
    >
      {state.items.length ? (
        state.items.map(({ value, label }, index) => (
          <li
            key={value}
            className={classNames(styles.dropdownMenuItem, {
              [styles.isFocused]: state.currentIndex === index,
              [styles.isSelected]: state.selection === value,
            })}
            {...getMenuItemProps({ value, index })}
          >
            {label}
          </li>
        ))
      ) : showCreateOption ? (
        <CreatableListItem />
      ) : (
        <li />
      )}
    </ul>
  );
};

export const SelectMenu = ({
  noOfVisibleItems = 9,
  itemSize = 28,
}: SelectMenuProps) => {
  const { state, getMenuProps, getMenuItemProps } = useSelect();
  const maxMenuHeight = itemSize * noOfVisibleItems;

  return (
    <ul
      className={styles.dropdownMenu}
      {...getMenuProps({ maxHeight: maxMenuHeight })}
    >
      {state.items.length ? (
        state.items.map(({ value, label }, index) => (
          <li
            key={value}
            className={classNames(styles.dropdownMenuItem, {
              [styles.isFocused]: state.currentIndex === index,
              [styles.isSelected]: state.selection === value,
            })}
            {...getMenuItemProps({ value, index })}
          >
            {label}
          </li>
        ))
      ) : (
        <li className={styles.dropdownMenuItem}>No options available</li>
      )}
    </ul>
  );
};
