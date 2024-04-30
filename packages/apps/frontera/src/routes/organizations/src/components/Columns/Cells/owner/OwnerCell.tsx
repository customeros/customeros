import { useRef, useMemo, useState, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { User } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetUsersQuery } from '@shared/graphql/getUsers.generated';
import { Select, getContainerClassNames } from '@ui/form/Select/Select';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
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

  const { data, isSuccess } = useGetUsersQuery(
    client,
    {
      pagination: {
        limit: 1000,
        page: 0,
      },
    },
    {
      enabled: !getUsers.hasFetched,
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

      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const targetOrgIndex = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.findIndex(
              (c) => c.metadata.id === id,
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
        queryClient.invalidateQueries({ queryKey });
      }, 1000);
      setIsEditing(false);
    },
  });

  const removeOrganizationOwner = useRemoveOrganizationOwnerMutation(client, {
    onMutate: () => {
      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            const pageIndex = getOrganization.pagination.page - 1;
            const owner = draft?.pages[
              pageIndex
            ]?.dashboardView_Organizations?.content.find(
              (c) => c.metadata.id === id,
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
        queryClient.invalidateQueries({ queryKey });
      }, 1000);

      setIsEditing(false);
    },
  });

  const handleSelect = useCallback(
    (option: SelectOption) => {
      if ((!option || !option.value) && !!prevSelection) {
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

  useEffect(() => {
    if (getUsers.hasFetched) return;
    setOrganizationsMeta(
      produce(organizationsMeta, (draft) => {
        draft.getUsers.hasFetched = true;
      }),
    );
  }, [isSuccess]);

  if (!isEditing) {
    return (
      <div className='flex w-full gap-1 items-center [&_.edit-button]:hover:opacity-100'>
        <p
          className={cn(
            value ? 'text-gray-700' : 'text-gray-400',
            'cursor-default',
          )}
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.label ?? 'Owner'}
        </p>
        <IconButton
          className='edit-button rounded-md opacity-0 '
          aria-label='erc'
          size='sm'
          variant='ghost'
          id='edit-button'
          onClick={() => setIsEditing(true)}
          icon={<Edit03 className='text-gray-500 size-3' />}
        />
      </div>
    );
  }

  return (
    <Select
      size='md'
      isClearable
      value={value}
      isLoading={
        setOrganizationOwner.isPending || removeOrganizationOwner.isPending
      }
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
      classNames={{
        container: () => getContainerClassNames('hover:border-transparent'),
      }}
    />
  );
};
