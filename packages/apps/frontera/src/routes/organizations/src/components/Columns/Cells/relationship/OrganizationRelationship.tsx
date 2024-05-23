import { useRef, useState, useEffect, useCallback } from 'react';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { Menu, MenuList, MenuItem, MenuButton } from '@ui/overlay/Menu/Menu';
import { useOrganizationQuery } from '@organization/graphql/organization.generated';
import { OrganizationRelationship as OrganizationRelationshipEnum } from '@graphql/types';
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
              content[index].relationship = payload.input.relationship;
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
    (option) => option.value === organization.relationship,
  );

  const handleSelect = useCallback(
    (option: SelectOption<OrganizationRelationshipEnum>) => {
      updateOrganization.mutate(
        OrganizationRowDTO.toUpdatePayload({
          ...organization,
          relationship: option.value,
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

  return (
    <div className='flex gap-1 items-center group'>
      <p
        className='cursor-default text-gray-700'
        onDoubleClick={() => setIsEditing(true)}
      >
        {value?.label || 'Unknown'}
      </p>
      <Menu open={isEditing} onOpenChange={setIsEditing}>
        <MenuButton asChild>
          <IconButton
            className={cn(
              'rounded-md opacity-0 group-hover:opacity-100',
              isEditing && 'opacity-100',
            )}
            aria-label='edit relationship'
            size='xxs'
            variant='ghost'
            id='edit-button'
            onClick={() => setIsEditing(true)}
            isLoading={updateOrganization.isPending}
            icon={<Edit03 className='text-gray-500' />}
          />
        </MenuButton>
        <MenuList>
          {relationshipOptions.map((option) => (
            <MenuItem
              key={option.value.toString()}
              onClick={() => handleSelect(option)}
            >
              {option.label}
            </MenuItem>
          ))}
        </MenuList>
      </Menu>
    </div>
  );
};
