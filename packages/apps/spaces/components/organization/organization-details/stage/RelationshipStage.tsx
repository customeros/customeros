import { useCallback } from 'react';
import {
  OrganizationRelationship,
  useSetStageToOrganizationRelationshipMutation,
  useRemoveStageFromOrganizationRelationshipMutation,
} from '@spaces/graphql';
import { stageOptions } from './util';
import {
  Select,
  SelectMenu,
  SelectInput,
  SelectWrapper,
} from '@spaces/ui-kit/select';

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
    (value: string) => {
      if (!relationship || !organizationId) return;

      if (!value) {
        handleRemoveStage();
      } else {
        handleAddStage(value);
      }
    },
    [handleRemoveStage, handleAddStage, relationship, organizationId],
  );

  return (
    <Select<string>
      options={stageOptions}
      onSelect={handleSelect}
      value={defaultValue ?? ''}
    >
      <SelectWrapper isHidden={!relationship}>
        <SelectInput
          placeholder='Stage'
          saving={loading}
          customStyles={{ marginLeft: 32 }}
        />
        <SelectMenu />
      </SelectWrapper>
    </Select>
  );
};
