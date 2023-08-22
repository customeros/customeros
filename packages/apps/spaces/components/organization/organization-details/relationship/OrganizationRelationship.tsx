import { useState, useCallback } from 'react';

import {
  OrganizationRelationship as Relationship,
  useAddRelationshipToOrganizationMutation,
  useRemoveOrganizationRelationshipMutation,
} from '@spaces/graphql';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';

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
  const [isEditing, setIsEditing] = useState(false);
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
            onCompleted: () => {
              setIsEditing(false);
            },
          });
        }

        addRelationshipToOrganization({
          variables: {
            organizationId,
            relationship,
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
          color={value ? 'gray.700' : 'gray.500'}
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.label ?? 'Relationship'}
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
      variant='unstyled'
      isLoading={loading}
      backspaceRemovesValue
      onChange={handleSelect}
      openMenuOnClick={false}
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
          boxSize: '3',
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
        menuList: (props) => ({
          ...props,
          w: '262px',
        }),
      }}
    />
  );
};
