import { useCallback } from 'react';
import {
  OrganizationRelationship,
  useSetStageToOrganizationRelationshipMutation,
  useRemoveStageFromOrganizationRelationshipMutation,
} from '@spaces/graphql';
import { stageOptions } from './util';

import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';

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
  const [setStageToRelationship, { loading }] =
    useSetStageToOrganizationRelationshipMutation();
  const [removeStageFromRelationship] =
    useRemoveStageFromOrganizationRelationshipMutation();
  const value = defaultValue
    ? stageOptions.find((o) => o.value === defaultValue)
    : null;

  const handleRemoveStage = useCallback(() => {
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

  const handleAddStage = useCallback(
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
    (option: SelectOption) => {
      if (!relationship || !organizationId) return;

      if (!option) {
        handleRemoveStage();
      } else {
        handleAddStage(option.value);
      }
    },
    [handleRemoveStage, handleAddStage, relationship, organizationId],
  );

  return (
    <Select
      size='sm'
      isClearable
      value={value}
      isLoading={loading}
      variant='unstyled'
      placeholder='Stage'
      backspaceRemovesValue
      options={stageOptions}
      onChange={handleSelect}
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
