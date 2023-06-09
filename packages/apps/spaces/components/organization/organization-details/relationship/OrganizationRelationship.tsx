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
import { SelectWrapper } from '@spaces/atoms/select/SelectWrapper';
import { SelectInput } from '@spaces/atoms/select/SelectInput';
import styles from '@spaces/atoms/select/select.module.scss';

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

const OrganizationSelectInput = () => {
  const { state } = useSelect();

  return (
    <>
      <SelectMenuItemIcon
        width='24'
        height='24'
        viewBox='0 0 24 24'
        name={state.selection as Relationship}
      />
      <SelectInput
        placeholder='Relationship'
        customStyles={{ paddingLeft: '8px' }}
      />
    </>
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

  const removeRelationship = useCallback(() => {
    removeOrganizationRelationship({
      variables: {
        organizationId,
        relationship: prevSelection,
      },
      update: (cache) => {
        const normalizedId = cache.identify({
          id: organizationId,
          __typename: 'Organization',
        });

        cache.modify({
          id: normalizedId,
          fields: {
            relationshipStages() {
              return [];
            },
          },
        });
        cache.gc();
      },
    });
  }, [removeOrganizationRelationship, organizationId, prevSelection]);

  const addRelationship = useCallback(
    (relationship: Relationship) => {
      if (relationship && relationship !== prevSelection) {
        if (prevSelection) {
          removeOrganizationRelationship({
            variables: {
              organizationId,
              relationship: prevSelection,
            },
          });
        }

        addRelationshipToOrganization({
          variables: {
            organizationId,
            relationship,
          },
          update: (cache) => {
            const normalizedId = cache.identify({
              id: organizationId,
              __typename: 'Organization',
            });

            cache.modify({
              id: normalizedId,
              fields: {
                relationshipStages() {
                  return [
                    {
                      __typename: 'OrganizationRelationshipStage',
                      relationship,
                      stage: null,
                    },
                  ];
                },
              },
            });
            cache.gc();
          },
        });
      }
    },
    [
      removeOrganizationRelationship,
      addRelationshipToOrganization,
      organizationId,
      prevSelection,
    ],
  );

  const handleSelect = useCallback(
    (relationship: Relationship) => {
      if (!relationship && prevSelection) {
        removeRelationship();
      } else {
        addRelationship(relationship);
        setPrevSelection(relationship);
      }
    },
    [prevSelection, addRelationship, removeRelationship],
  );

  return (
    <Select<Relationship>
      onSelect={handleSelect}
      value={defaultValue}
      options={relationshipOptions}
    >
      <SelectWrapper>
        <OrganizationSelectInput />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
