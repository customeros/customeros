import { useRef } from 'react';

import { produce } from 'immer';
import merge from 'lodash/merge';
import { useWillUnmount } from 'rooks';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { useAddSocialMutation } from '@organization/graphql/addOrganizationSocial.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/graphql/organization.generated';
import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '@organizations/graphql/getOrganizations.generated';

interface UseAboutPanelMethodsOptions {
  id: string;
}

export const useAboutPanelMethods = ({ id }: UseAboutPanelMethodsOptions) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();
  const timeout = useRef<NodeJS.Timeout>();

  const queryKey = useOrganizationQuery.getKey({ id });
  const organizationsQueryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );

  const invalidateQuery = () => queryClient.invalidateQueries({ queryKey });

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
              (o) => o.metadata.id === id,
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
      if (timeout.current) clearTimeout(timeout.current);

      timeout.current = setTimeout(() => {
        invalidateQuery();
        queryClient.invalidateQueries({ queryKey: organizationsQueryKey });
      }, 1000);
    },
  });

  const addSocial = useAddSocialMutation(client, {
    onSuccess: invalidateQuery,
  });

  useWillUnmount(() => {
    if (timeout.current) clearTimeout(timeout.current);
  });

  return { updateOrganization, addSocial, invalidateQuery };
};
