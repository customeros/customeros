import { useCallback, useState } from 'react';
import {
  OrganizationRelationship,
  useSetStageToOrganizationRelationshipMutation,
  useRemoveStageFromOrganizationRelationshipMutation,
} from '@spaces/graphql';

import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';

import { customerStageOptions, otherStageOptions } from './util';

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
  const [isEditing, setIsEditing] = useState(false);
  const [setStageToRelationship, { loading }] =
    useSetStageToOrganizationRelationshipMutation();
  const [removeStageFromRelationship] =
    useRemoveStageFromOrganizationRelationshipMutation();
  const options =
    relationship === 'CUSTOMER' ? customerStageOptions : otherStageOptions;
  const value = defaultValue
    ? options.find((o) => o.value === defaultValue)
    : null;

  const handleRemoveStage = useCallback(() => {
    if (!relationship || !organizationId) return;

    removeStageFromRelationship({
      variables: {
        organizationId,
        relationship: relationship as OrganizationRelationship,
      },
      onCompleted: () => {
        setIsEditing(false);
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
        onCompleted: () => {
          setIsEditing(false);
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

  if (!isEditing) {
    return (
      <Flex
        w='full'
        gap='1'
        align='center'
        _hover={{
          '& #edit-button': {
            opacity: 1,
          },
        }}
      >
        <Text
          cursor='default'
          color={value ? 'gray.700' : 'gray.400'}
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.label ?? 'Stage'}
        </Text>
        <IconButton
          aria-label='erc'
          size='xs'
          borderRadius='md'
          minW='4'
          w='4'
          minH='4'
          h='4'
          opacity='0'
          variant='ghost'
          id='edit-button'
          onClick={() => setIsEditing(true)}
          icon={<Icons.Edit3 color='gray.500' boxSize='3' />}
        />
      </Flex>
    );
  }

  return (
    <Select
      size='sm'
      isClearable
      value={value}
      autoFocus
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          setIsEditing(false);
        }
      }}
      defaultMenuIsOpen
      onBlur={() => setIsEditing(false)}
      isLoading={loading}
      variant='unstyled'
      placeholder='Stage'
      backspaceRemovesValue
      options={options}
      onChange={handleSelect}
      openMenuOnClick={false}
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
          boxSize: '3',
        }),
        placeholder: (props) => ({
          ...props,
          ml: 0,
          color: 'gray.400',
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
