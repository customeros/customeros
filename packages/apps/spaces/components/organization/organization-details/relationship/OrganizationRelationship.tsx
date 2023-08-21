import { useState, useCallback } from 'react';

import {
  OrganizationRelationship as Relationship,
  useAddRelationshipToOrganizationMutation,
  useRemoveOrganizationRelationshipMutation,
} from '@spaces/graphql';
import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';

import { relationshipOptions } from './util';

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
  const [addRelationshipToOrganization, { loading }] =
    useAddRelationshipToOrganizationMutation();
  const [removeOrganizationRelationship] =
    useRemoveOrganizationRelationshipMutation();
  const value = defaultValue
    ? relationshipOptions.find((o) => o.value === defaultValue)
    : null;

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
    (option: SelectOption<Relationship>) => {
      if (!option && prevSelection) {
        removeRelationship();
      } else {
        addRelationship(option.value);
        setPrevSelection(option.value);
      }
    },
    [prevSelection, addRelationship, removeRelationship],
  );

  return (
    <Select
      size='sm'
      isClearable
      value={value}
      variant='unstyled'
      isLoading={loading}
      backspaceRemovesValue
      onChange={handleSelect}
      placeholder='Relationship'
      options={relationshipOptions}
      chakraStyles={{
        valueContainer: (props) => ({
          ...props,
          p: 0,
        }),
        singleValue: (props) => ({
          ...props,
          paddingBottom: 0,
          ml: 0,
        }),
        control: (props) => ({
          ...props,
          minH: '0',
        }),
        clearIndicator: (props) => ({
          ...props,
          display: 'none',
        }),
        placeholder: (props) => ({
          ...props,
          ml: 0,
          color: 'gray.500',
        }),
        inputContainer: (props) => ({
          ...props,
          py: 0,
          ml: 0,
        }),
      }}
    />
  );
};
