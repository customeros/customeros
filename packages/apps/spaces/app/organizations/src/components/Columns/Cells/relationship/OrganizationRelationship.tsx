import { useRef, useState, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Edit03 } from '@ui/media/icons/Edit03';
import { Select } from '@ui/form/Select/Select';
import { getContainerClassNames } from '@ui/form/Select';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { getMenuListClassNames } from '@ui/form/MultiCreatableSelect/MultiCreatableSelect2';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import {
  OrganizationRowDTO,
  GetOrganizationRowResult,
} from '@organizations/util/Organization.dto';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

import { relationshipOptions } from './util';

interface OrganizationRelationshipProps {
  organization: GetOrganizationRowResult;
}

export const OrganizationRelationship = ({
  organization,
}: OrganizationRelationshipProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [organizationsMeta] = useOrganizationsMeta();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const updateOrganization = useUpdateOrganizationMutation(client, {
    onMutate: (payload) => {
      queryClient.cancelQueries({ queryKey });

      const previousOrganizations =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          const pageIndex =
            organizationsMeta.getOrganization.pagination.page - 1;

          return produce(old, (draft) => {
            const content =
              draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content;
            const index = content?.findIndex(
              (item) => item.metadata.id === payload.input.id,
            );

            if (content && index !== undefined && index > -1) {
              content[index].isCustomer = payload.input.isCustomer;
            }
          });
        },
      );

      return { previousOrganizations };
    },
    onError: (_, __, context) => {
      if (context?.previousOrganizations) {
        queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
          queryKey,
          context.previousOrganizations,
        );
      }
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({
          queryKey: useOrganizationQuery.getKey({
            id: organization.metadata.id,
          }),
        });
      }, 500);
    },
  });

  const { getOrganization } = organizationsMeta;
  const queryKey = useInfiniteGetOrganizationsQuery.getKey(getOrganization);

  const value = relationshipOptions.find(
    (option) => option.value === organization.isCustomer,
  );

  const handleSelect = useCallback(
    (option: SelectOption<boolean>) => {
      updateOrganization.mutate(
        OrganizationRowDTO.toUpdatePayload({
          ...organization,
          isCustomer: option.value,
        }),
      );
    },
    [updateOrganization, organization],
  );

  useEffect(() => {
    return () => {
      timeoutRef.current && clearTimeout(timeoutRef.current);
    };
  }, []);

  if (!isEditing) {
    return (
      <div className='flex w-full gap-1 items-center [&_.edit-button]:hover:opacity-100'>
        <p
          className='cursor-default text-gray-700 group'
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.value ? 'Customer' : 'Prospect'}
        </p>
        <IconButton
          className='edit-button rounded-md opacity-0'
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
      size='sm'
      isClearable
      defaultValue={value}
      autoFocus
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          setIsEditing(false);
        }
      }}
      defaultMenuIsOpen
      onBlur={() => setIsEditing(false)}
      isLoading={updateOrganization.isPending}
      backspaceRemovesValue
      onChange={handleSelect}
      openMenuOnClick={false}
      placeholder='Relationship'
      options={relationshipOptions}
      classNames={{
        container: () => getContainerClassNames('hover:border-transparent'),
        menuList: () => getMenuListClassNames('w-[262px]'),
      }}
    />
  );
};
