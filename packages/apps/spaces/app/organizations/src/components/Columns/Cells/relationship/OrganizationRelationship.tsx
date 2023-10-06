import { useState, useCallback, useRef, useEffect } from 'react';
import { useQueryClient, InfiniteData } from '@tanstack/react-query';
import { produce } from 'immer';

import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { OrganizationRelationship as Relationship } from '@graphql/types';
import { useAddRelationshipMutation } from '@organization/src/graphql/addRelationship.generated';
import { useRemoveRelationshipMutation } from '@organization/src/graphql/removeRelationship.generated';
import {
  useInfiniteGetOrganizationsQuery,
  GetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';

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
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();
  const addTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const removeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const [isEditing, setIsEditing] = useState(false);
  const [prevSelection, setPrevSelection] =
    useState<Relationship>(defaultValue);

  const { getOrganization } = organizationsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getOrganization);

  const addRelationship = useAddRelationshipMutation(client, {
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

            relationships?.splice(0, 1, {
              __typename: 'OrganizationRelationshipStage',
              relationship: payload.relationship,
              stage: null,
            });
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
  const removeRelationship = useRemoveRelationshipMutation(client, {
    onMutate: () => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries(queryKey);
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            draft?.pages[pageIndex]?.dashboardView_Organizations?.content
              .find((c) => c.id === organizationId)
              ?.relationshipStages?.splice(0, 1);
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

      setIsEditing(false);
    },
  });

  const value = defaultValue
    ? relationshipOptions.find((o) => o.value === defaultValue)
    : null;

  const add = useCallback(
    (relationship: Relationship) => {
      if (relationship && relationship !== prevSelection) {
        if (prevSelection) {
          removeRelationship.mutate({
            organizationId,
            relationship: prevSelection,
          });
        }

        addRelationship.mutate({
          organizationId,
          relationship,
        });
      }
    },
    [addRelationship, organizationId, prevSelection],
  );

  const handleSelect = useCallback(
    (option: SelectOption<Relationship>) => {
      if (!option && prevSelection) {
        removeRelationship.mutate({
          organizationId,
          relationship: prevSelection,
        });
      } else {
        add(option.value);
        setPrevSelection(option.value);
      }
    },
    [prevSelection, add, removeRelationship],
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
      isLoading={addRelationship.isLoading}
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
          color: 'gray.400',
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
