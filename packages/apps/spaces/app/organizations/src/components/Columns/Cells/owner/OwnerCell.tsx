import { useRef, useMemo, useState, useCallback } from 'react';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { User } from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Select } from '@ui/form/SyncSelect/Select';
import { SelectOption } from '@shared/types/SelectOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useGetUsersQuery } from '@organizations/graphql/getUsers.generated';
import { useSetOrganizationOwnerMutation } from '@organizations/graphql/setOrganizationOwner.generated';
import { useRemoveOrganizationOwnerMutation } from '@organizations/graphql/removeOrganizationOwner.generated';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
}

export const OwnerCell = ({ id, owner }: OwnerProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const addTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const removeTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();

  const [isEditing, setIsEditing] = useState(false);
  const [prevSelection, setPrevSelection] = useState(owner);

  const { getOrganization, getUsers } = organizationsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getOrganization);

  const { data } = useGetUsersQuery(
    client,
    {
      pagination: {
        limit: 1000,
        page: 1,
      },
    },
    {
      enabled: !getUsers.hasFetched,
      onSuccess: () => {
        if (getUsers.hasFetched) return;
        setOrganizationsMeta(
          produce(organizationsMeta, (draft) => {
            draft.getUsers.hasFetched = true;
          }),
        );
      },
    },
  );

  const options = useMemo(() => {
    return data?.users?.content
      ?.filter(
        (e) => Boolean(e.firstName) || Boolean(e.lastName) || Boolean(e.name),
      )
      ?.map((o) => ({
        value: o.id,
        label: `${o.name ?? o.firstName + ' ' + o.lastName}`.trim(),
      }))
      ?.sort((a, b) => a.label.localeCompare(b.label));
  }, [data]);

  const value = owner ? options?.find((o) => o.value === owner.id) : null;

  const setOrganizationOwner = useSetOrganizationOwnerMutation(client, {
    onMutate: (payload) => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries(queryKey);
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const targetOrgIndex = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.findIndex(
              (c) => c.id === id,
            );
            if (typeof targetOrgIndex === 'undefined' || targetOrgIndex < 0)
              return;

            const targetOrg =
              draft?.pages[pageIndex]?.dashboardView_Organizations?.content[
                targetOrgIndex
              ];
            if (!targetOrg) return;

            const foundOption = options?.find(
              (o) => o.value === payload.userId,
            );

            if (!foundOption) return;
            const [firstName, lastName] = foundOption.label.split(' ');

            const ownerItem = targetOrg?.owner;
            if (!ownerItem) {
              targetOrg.owner = {
                id: payload.userId,
                firstName,
                lastName,
              };
            } else {
              ownerItem.id = payload.userId;
              ownerItem.firstName = firstName;
              ownerItem.lastName = lastName;
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

  const removeOrganizationOwner = useRemoveOrganizationOwnerMutation(client, {
    onMutate: () => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries(queryKey);
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const owner = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.find(
              (c) => c.id === id,
            )?.owner;

            if (owner) {
              owner.id = '';
              owner.firstName = '';
              owner.lastName = '';
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

      setIsEditing(false);
    },
  });

  const handleSelect = useCallback(
    (option: SelectOption) => {
      if (!option && !!prevSelection) {
        removeOrganizationOwner.mutate({
          organizationId: id,
        });
        setPrevSelection(null);
      } else {
        setOrganizationOwner.mutate({
          userId: option.value,
          organizationId: id,
        });
        setPrevSelection(owner);
      }
    },
    [prevSelection, owner, removeOrganizationOwner, setOrganizationOwner, id],
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
          {value?.label ?? 'Owner'}
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
      isLoading={
        setOrganizationOwner.isLoading || removeOrganizationOwner.isLoading
      }
      variant='unstyled'
      placeholder='Owner'
      autoFocus
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          setIsEditing(false);
        }
      }}
      defaultMenuIsOpen
      onBlur={() => setIsEditing(false)}
      backspaceRemovesValue
      openMenuOnClick={false}
      onChange={handleSelect}
      options={options}
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
