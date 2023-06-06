import type { PropsWithChildren } from 'react';
import classNames from 'classnames';

import { Select, useSelect } from '@spaces/atoms/select';

import { relationshipOptions } from './util';
import styles from './organization-relationship.module.scss';

interface SelectMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
}

const SelectMenu = ({
  noOfVisibleItems = 7,
  itemSize = 25,
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
        <li className={styles.dropdownMenuItem} data-dropdown='menuitem'>
          No options available
        </li>
      )}
    </ul>
  );
};

const SelectInput = () => {
  const { state, getInputProps, autofillValue } = useSelect();

  return (
    <>
      <span
        role='textbox'
        placeholder='Relationship'
        contentEditable={state.isEditing}
        className={classNames(styles.dropdownInput)}
        {...getInputProps()}
      />
      <span className={styles.autofill}>{autofillValue}</span>
    </>
  );
};

const SelectWrapper = ({ children }: PropsWithChildren) => {
  const { getWrapperProps } = useSelect();

  return (
    <div {...getWrapperProps()} className={styles.dropdownWrapper}>
      {children}
    </div>
  );
};

interface OrganizationRelationshipProps {
  defaultValue?: string;
}

export const OrganizationRelationship = ({
  defaultValue,
}: OrganizationRelationshipProps) => {
  return (
    <Select defaultValue={defaultValue} options={relationshipOptions}>
      <SelectWrapper>
        <SelectInput />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
