import merge from 'lodash/merge';
import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/src/graphql/organization.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { useAddSocialMutation } from '@organization/src/graphql/addOrganizationSocial.generated';

interface UseAboutPanelMethodsOptions {
  id: string;
}

export const useAboutPanelMethods = ({ id }: UseAboutPanelMethodsOptions) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();

  const queryKey = useOrganizationQuery.getKey({ id });
  const organizationsQueryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );

  const invalidateQuery = () => queryClient.invalidateQueries(queryKey);

  const updateOrganization = useUpdateOrganizationMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      const previousEntries =
        queryClient.getQueryData<OrganizationQuery>(queryKey);

      const previousOrganizationsEntries = queryClient.getQueryData<
        InfiniteData<GetOrganizationsQuery>
      >(organizationsQueryKey);

      queryClient.setQueryData<OrganizationQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          merge(draft, input);
        });
      });
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        organizationsQueryKey,
        (currentCache) => {
          return produce(currentCache, (draft) => {
            const pageIndex =
              organizationsMeta.getOrganization.pagination.page - 1;
            const foundIndex = draft?.pages?.[
              pageIndex
            ]?.dashboardView_Organizations?.content?.findIndex(
              (o) => o.id === id,
            );

            if (typeof foundIndex === 'undefined' || foundIndex < 0) return;
            merge(
              draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content?.[
                foundIndex
              ],
              input,
            );
          });
        },
      );

      return { previousEntries, previousOrganizationsEntries };
    },

    onError: (_, __, context) => {
      queryClient.setQueryData(queryKey, context?.previousEntries);
      queryClient.setQueryData(
        organizationsQueryKey,
        context?.previousOrganizationsEntries,
      );
    },
    onSettled: () => {
      invalidateQuery();
      queryClient.invalidateQueries(organizationsQueryKey);
    },
  });

  const addSocial = useAddSocialMutation(client, {
    onSuccess: invalidateQuery,
  });

  return { updateOrganization, addSocial, invalidateQuery };
};
