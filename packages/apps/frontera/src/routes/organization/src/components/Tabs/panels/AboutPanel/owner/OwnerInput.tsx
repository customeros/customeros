import { useRef, useMemo, useCallback } from 'react';

import { useQueryClient } from '@tanstack/react-query';

import { User } from '@graphql/types';
import { Select } from '@ui/form/Select';
import { User02 } from '@ui/media/icons/User02';
import { SelectOption } from '@shared/types/SelectOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetUsersQuery } from '@shared/graphql/getUsers.generated';
import { useSetOrganizationOwnerMutation } from '@organizations/graphql/setOrganizationOwner.generated';
import { useRemoveOrganizationOwnerMutation } from '@organizations/graphql/removeOrganizationOwner.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/graphql/organization.generated';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
  invalidateQuery: () => void;
}

export const OwnerInput = ({ id, owner, invalidateQuery }: OwnerProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const addTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const removeTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryKey = useOrganizationQuery.getKey({ id });
  const { data } = useGetUsersQuery(client, {
    pagination: {
      limit: 1000,
      page: 1,
    },
  });

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
      const newOwner = data?.users?.content?.find(
        (o) => o.id === payload.userId,
      );
      const organization =
        queryClient.getQueryData<OrganizationQuery>(queryKey);
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<OrganizationQuery>(queryKey, (oldData) => {
        if (!oldData || !oldData?.organization) return;

        return {
          ...oldData,
          organization: {
            ...(oldData?.organization ?? {}),
            owner: newOwner,
          },
        };
      });

      return { organization };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<OrganizationQuery>(
        queryKey,
        () => context?.organization,
      );
    },
    onSettled: () => {
      if (addTimeoutRef.current) {
        clearTimeout(addTimeoutRef.current);
      }
      addTimeoutRef.current = setTimeout(() => {
        invalidateQuery();
      }, 1000);
    },
  });

  const removeOrganizationOwner = useRemoveOrganizationOwnerMutation(client, {
    onMutate: () => {
      const organization =
        queryClient.getQueryData<OrganizationQuery>(queryKey);
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<OrganizationQuery>(queryKey, (oldData) => {
        if (!oldData || !oldData?.organization) return;

        return {
          ...oldData,
          organization: {
            ...(oldData?.organization ?? {}),
            owner: null,
          },
        };
      });

      return { organization };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<OrganizationQuery>(
        queryKey,
        () => context?.organization,
      );
    },
    onSettled: () => {
      if (removeTimeoutRef.current) {
        clearTimeout(removeTimeoutRef.current);
      }
      removeTimeoutRef.current = setTimeout(() => {
        invalidateQuery();
      }, 1000);
    },
  });

  const handleSelect = useCallback(
    (option: SelectOption) => {
      if (!option || !option.value) {
        removeOrganizationOwner.mutate({
          organizationId: id,
        });
      } else {
        setOrganizationOwner.mutate({
          userId: option.value,
          organizationId: id,
        });
      }
    },
    [owner, removeOrganizationOwner, setOrganizationOwner, id],
  );

  return (
    <Select
      isClearable
      value={value}
      isLoading={setOrganizationOwner.isPending}
      placeholder='Owner'
      backspaceRemovesValue
      onChange={handleSelect}
      options={options}
      leftElement={<User02 className='text-gray-500 mr-3' />}
    />
  );
};
