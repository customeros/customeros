'use client';

import { useParams } from 'next/navigation';

import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Select } from '@ui/form/SyncSelect';
import { Spinner } from '@ui/feedback/Spinner';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { GetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { relationshipOptions } from '@organizations/components/Columns/Cells/relationship/util';
import { contractButtonSelect } from '@organization/src/components/Tabs/shared/contractSelectStyles';

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
              (item) => item.id === payload.input.id,
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

  return (
    <Box>
      <Select
        isSearchable={false}
        isClearable={false}
        isMulti={false}
        value={selectedValue}
        onChange={(value) =>
          updateOrganization.mutate({
            input: {
              id,
              isCustomer: value?.value,
              patch: true,
            },
          })
        }
        options={relationshipOptions}
        chakraStyles={{
          ...contractButtonSelect,
          container: (props, state) => {
            const isCustomer = state.getValue()[0]?.value;

            return {
              ...props,
              px: 2,
              py: '1px',
              border: '1px solid',
              borderColor: isCustomer ? 'success.200' : 'gray.300',
              backgroundColor: isCustomer ? 'success.50' : 'transparent',
              color: isCustomer ? 'success.700' : 'gray.500',

              borderRadius: '2xl',
              fontSize: 'xs',
              maxHeight: '22px',

              '& > div': {
                p: 0,
                border: 'none',
                fontSize: 'xs',
                maxHeight: '22px',
                minH: 'auto',
              },
            };
          },
          valueContainer: (props, state) => {
            const isCustomer = state.getValue()[0]?.value;

            return {
              ...props,
              p: 0,
              border: 'none',
              fontSize: 'xs',
              maxHeight: '22px',
              minH: 'auto',
              color: isCustomer ? 'success.700' : 'gray.500',
            };
          },
          singleValue: (props) => {
            return {
              ...props,
              maxHeight: '22px',
              p: 0,
              minH: 'auto',
              color: 'inherit',
            };
          },
          menuList: (props) => {
            return {
              ...props,
              w: 'fit-content',
              left: '-32px',
            };
          },
        }}
        leftElement={
          updateOrganization.isPending ? (
            <Spinner
              size='xs'
              mr={1}
              color={selectedValue?.value ? 'success.500' : 'gray.400'}
            />
          ) : selectedValue?.value ? (
            <ActivityHeart color='success.500' mr='1' />
          ) : null
        }
      />
    </Box>
  );
};
