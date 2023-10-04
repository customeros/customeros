import { useCallback, useState, useRef, useEffect } from 'react';
import { useQueryClient, InfiniteData } from '@tanstack/react-query';
import { produce } from 'immer';

import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';
import { OrganizationRelationship } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useSetRelationshipStageMutation } from '@organization/src/graphql/setRelationshipStage.generated';
import { useRemoveRelationshipStageMutation } from '@organization/src/graphql/removeRelationshipStage.generated';
import {
  useInfiniteGetOrganizationsQuery,
  GetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';
import { useOrganizationsMeta } from '@organizations/shared/state';

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
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();
  const addTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const removeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const [isEditing, setIsEditing] = useState(false);

  const { getOrganization } = organizationsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getOrganization);

  const addStage = useSetRelationshipStageMutation(client, {
    onMutate: (payload) => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries(queryKey);
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const relationships = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.find(
              (c) => c.id === organizationId,
            )?.relationshipStages;

            const relationship = relationships?.[0];
            if (relationship) {
              relationship.stage = payload.stage;
            }
          });
        },
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        context?.previousEntries,
      );
    },
    onSettled: () => {
      if (addTimeoutRef.current) {
        clearTimeout(addTimeoutRef.current);
      }
      addTimeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
      }, 1000);
      setIsEditing(false);
    },
  });
  const removeStage = useRemoveRelationshipStageMutation(client, {
    onMutate: () => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries(queryKey);
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const relationships = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.find(
              (c) => c.id === organizationId,
            )?.relationshipStages;

            const relationship = relationships?.[0];
            if (relationship) {
              relationship.stage = null;
            }
          });
        },
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        context?.previousEntries,
      );
    },
    onSettled: () => {
      if (removeTimeoutRef.current) {
        clearTimeout(removeTimeoutRef.current);
      }
      removeTimeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
      }, 1000);
      ``;
      setIsEditing(false);
    },
  });

  const options =
    relationship === 'CUSTOMER' ? customerStageOptions : otherStageOptions;
  const value = defaultValue
    ? options.find((o) => o.value === defaultValue)
    : null;

  const handleRemoveStage = useCallback(() => {
    if (!relationship || !organizationId) return;

    removeStage.mutate({
      organizationId,
      relationship: relationship as OrganizationRelationship,
    });
  }, [organizationId, relationship, removeStage, queryKey, queryClient]);

  const handleAddStage = useCallback(
    (value: string) => {
      if (!relationship || !organizationId) return;

      addStage.mutate({
        organizationId,
        relationship: relationship as OrganizationRelationship,
        stage: value,
      });
    },
    [organizationId, relationship, addStage, queryKey, queryClient],
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

  useEffect(() => {
    return () => {
      addTimeoutRef.current && clearTimeout(addTimeoutRef.current);
      removeTimeoutRef.current && clearTimeout(removeTimeoutRef.current);
    };
  }, []);

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
      isLoading={addStage.isLoading}
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
