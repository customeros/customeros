import { PropsWithChildren } from 'react';
import { useState, useCallback } from 'react';
import classNames from 'classnames';

import { Select, useSelect } from '@spaces/atoms/select';

import {
  OrganizationRelationship as Relationship,
  useAddRelationshipToOrganizationMutation,
  useRemoveOrganizationRelationshipMutation,
} from '@spaces/graphql';

import { relationshipOptions } from './util';
import { SelectMenuItemIcon } from './SelectMenuItemIcon';
import styles from './organization-relationship.module.scss';

interface SelectMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
}

const SelectMenu = ({
  noOfVisibleItems = 7,
  itemSize = 38,
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
            <SelectMenuItemIcon
              width='16'
              height='16'
              viewBox='0 0 24 24'
              name={value as Relationship}
            />{' '}
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
      <SelectMenuItemIcon
        width='16'
        height='16'
        viewBox='0 0 24 24'
        name={state.selection as Relationship}
      />
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
  defaultValue: Relationship;
  organizationId: string;
}

export const OrganizationRelationship = ({
  defaultValue,
  organizationId,
}: OrganizationRelationshipProps) => {
  const [prevSelection, setPrevSelection] =
    useState<Relationship>(defaultValue);
  const [addRelationshipToOrganization] =
    useAddRelationshipToOrganizationMutation();
  const [removeOrganizationRelationship] =
    useRemoveOrganizationRelationshipMutation();

  const handleSelect = useCallback(
    (relationship: Relationship) => {
      if (!relationship && prevSelection) {
        removeOrganizationRelationship({
          variables: {
            organizationId,
            relationship: prevSelection,
          },
        });
      } else {
        if (relationship && relationship !== prevSelection) {
          addRelationshipToOrganization({
            variables: {
              organizationId,
              relationship,
            },
          });
        }
        setPrevSelection(relationship);
      }
    },
    [
      prevSelection,
      addRelationshipToOrganization,
      removeOrganizationRelationship,
      organizationId,
    ],
  );

  return (
    <Select<Relationship>
      onSelect={handleSelect}
      defaultValue={defaultValue}
      options={relationshipOptions}
    >
      <SelectWrapper>
        <SelectInput />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
