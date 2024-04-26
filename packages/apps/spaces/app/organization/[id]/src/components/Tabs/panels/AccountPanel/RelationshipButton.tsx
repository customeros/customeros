'use client';

import { useParams } from 'next/navigation';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { GetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { relationshipOptions } from '@organizations/components/Columns/Cells/relationship/util';

export const RelationshipButton = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const id = useParams()?.id as string;
  const { data } = useOrganizationQuery(client, { id });
  const queryKey = useOrganizationQuery.getKey({ id });
  const [organizationsMeta] = useOrganizationsMeta();

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
    onSettled: () =>
      queryClient.invalidateQueries({
        queryKey: useOrganizationQuery.getKey({ id }),
      }),
  });
  const selectedValue = relationshipOptions.find(
    (option) => option.value === data?.organization?.isCustomer,
  );
  const spinnerColors = selectedValue?.value
    ? 'text-success-500 fill-succes-700'
    : 'text-gray-400 fill-gray-700';

  return (
    <div>
      <Menu>
        <MenuButton asChild>
          <Button
            variant='outline'
            size='xxs'
            colorScheme={
              selectedValue?.label === 'Customer' ? 'success' : 'gray'
            }
            className={cn(
              selectedValue?.label === 'Customer'
                ? 'text-success-500'
                : 'text-gray-500',
              'rounded-full font-normal  text-ellipsis mb-[2.5px]',
            )}
            leftIcon={
              updateOrganization.isPending ? (
                <Spinner
                  label='Organization loading'
                  size='sm'
                  className={cn(spinnerColors)}
                />
              ) : selectedValue?.value ? (
                <ActivityHeart className='text-success-500' />
              ) : (
                <></>
              )
            }
          >
            {selectedValue?.label}
          </Button>
        </MenuButton>
        <MenuList className='p-2'>
          {relationshipOptions.map((option, idx) => (
            <MenuItem
              className={cn(
                selectedValue?.label === option.label
                  ? 'text-primary-600 bg-primary-50 hover:bg-primary-50 '
                  : 'hover:bg-gray-100',
                'px-2 py-1 border border-transparent hover:border-gray-200 hover:border rounded-md',
              )}
              key={idx}
              onClick={() =>
                updateOrganization.mutate({
                  input: {
                    id,
                    isCustomer: option.value,
                    patch: true,
                  },
                })
              }
            >
              {option.label}
            </MenuItem>
          ))}
        </MenuList>
      </Menu>
    </div>
  );
};

// borderColor: isCustomer ? 'success.200' : 'gray.300',
// backgroundColor: isCustomer ? 'success.50' : 'transparent',
// color: isCustomer ? 'success.700' : 'gray.500',
