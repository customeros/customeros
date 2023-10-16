import { useState, useCallback, useRef, useEffect } from 'react';
import { useQueryClient, InfiniteData } from '@tanstack/react-query';
import { produce } from 'immer';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Select } from '@ui/form/SyncSelect/Select';

import {
  OrganizationRowDTO,
  GetOrganizationRowResult,
} from '@organizations/util/Organization.dto';
import { SelectOption } from '@shared/types/SelectOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';

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
      console.log('onMutate');
      queryClient.cancelQueries(queryKey);

      const previousOrganizations =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          const pageIndex = organizationsMeta.getOrganization.pagination.page;
          return produce(old, (draft) => {
            const content =
              draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content;
            const index = content?.findIndex(
              (item) => item.id === payload.input.id,
            );

            if (content && index) {
              content[index].isCustomer = payload.input.isCustomer;
            }
          });
        },
      );

      return { previousOrganizations };
    },
    onError: (_, __, context) => {
      console.log('onError');
      if (context?.previousOrganizations) {
        queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
          queryKey,
          context.previousOrganizations,
        );
      }
    },
    onSettled: () => {
      console.log('onSettled');
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
        queryClient.invalidateQueries(
          useOrganizationQuery.getKey({ id: organization.id }),
        );
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
      console.log('handleSelect');
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
          color='gray.700'
          onDoubleClick={() => setIsEditing(true)}
        >
          {value?.value ? 'Customer' : 'Prospect'}
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
      defaultValue={value}
      autoFocus
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          setIsEditing(false);
        }
      }}
      defaultMenuIsOpen
      onBlur={() => setIsEditing(false)}
      variant='unstyled'
      isLoading={updateOrganization.isLoading}
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
