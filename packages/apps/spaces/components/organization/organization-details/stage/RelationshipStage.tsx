import type { PropsWithChildren } from 'react';
import { useCallback } from 'react';
import classNames from 'classnames';

import { Select, useSelect } from '@spaces/atoms/select';
import {
  OrganizationRelationship,
  useSetStageToOrganizationRelationshipMutation,
  useRemoveStageFromOrganizationRelationshipMutation,
} from '@spaces/graphql';

import { stageOptions } from './util';
import styles from './organization-stage.module.scss';

interface SelectMenuProps {
  noOfVisibleItems?: number;
  itemSize?: number;
}

const SelectMenu = ({
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
        placeholder='Stage'
        contentEditable={state.isEditing}
        className={classNames(styles.dropdownInput)}
        {...getInputProps()}
      />
      <span className={styles.autofill}>{autofillValue}</span>
    </>
  );
};

const SelectWrapper = ({
  children,
  isVisible,
}: PropsWithChildren<{ isVisible?: boolean }>) => {
  const { getWrapperProps } = useSelect();

  return (
    <div
      {...getWrapperProps()}
      className={styles.dropdownWrapper}
      style={{ visibility: isVisible ? 'visible' : 'hidden' }}
    >
      {children}
    </div>
  );
};

interface RelationshipStageProps {
  defaultValue?: string | null;
  relationship?: string;
  organizationId?: string;
}

export const RelationshipStage = ({
  defaultValue,
  relationship,
  organizationId,
}: RelationshipStageProps) => {
  const [setStageToRelationship] =
    useSetStageToOrganizationRelationshipMutation();
  const [removeStageFromRelationship] =
    useRemoveStageFromOrganizationRelationshipMutation();

  const removeStage = useCallback(() => {
    if (!relationship || !organizationId) return;

    removeStageFromRelationship({
      variables: {
        organizationId,
        relationship: relationship as OrganizationRelationship,
      },
      update: (cache) => {
        const normalizedId = cache.identify({
          id: organizationId,
          __typename: 'Organization',
        });

        cache.modify({
          id: normalizedId,
          fields: {
            relationshipStages(v) {
              return [
                {
                  ...v?.[0],
                  stage: null,
                },
              ];
            },
          },
        });
        cache.gc();
      },
    });
  }, [organizationId, relationship, removeStageFromRelationship]);

  const addStage = useCallback(
    (value: string) => {
      if (!relationship || !organizationId) return;

      setStageToRelationship({
        variables: {
          organizationId,
          relationship: relationship as OrganizationRelationship,
          stage: value,
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
                    relationship: relationship as OrganizationRelationship,
                    stage: value,
                  },
                ];
              },
            },
          });
          cache.gc();
        },
      });
    },
    [organizationId, relationship, setStageToRelationship],
  );

  const handleSelect = useCallback(
    (value: string) => {
      if (!relationship || !organizationId) return;

      if (!value) {
        removeStage();
      } else {
        addStage(value);
      }
    },
    [removeStage, addStage, relationship, organizationId],
  );

  return (
    <Select<string>
      options={stageOptions}
      onSelect={handleSelect}
      value={defaultValue ?? ''}
    >
      <SelectWrapper isVisible={!!relationship}>
        <SelectInput />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
