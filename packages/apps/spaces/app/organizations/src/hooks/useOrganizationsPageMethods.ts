import set from 'lodash/set';
import { produce } from 'immer';
import { useRouter } from 'next/navigation';
import { useQueryClient, InfiniteData } from '@tanstack/react-query';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '../graphql/getOrganizations.generated';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useCreateOrganizationMutation } from '../graphql/createOrganization.generated';
import { useHideOrganizationsMutation } from '../graphql/hideOrganizations.generated';
import { useMergeOrganizationsMutation } from '../graphql/mergeOrganizations.generated';

export const useOrganizationsPageMethods = () => {
  const { push } = useRouter();
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta, setOrganizationsMeta] = useOrganizationsMeta();

  const queryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );

  const createOrganization = useCreateOrganizationMutation(client, {
    onMutate: () => {
      const pageIndex = organizationsMeta.getOrganization.pagination.page - 1;
      queryClient.cancelQueries(queryKey);

      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            if (!draft) return;

            const page = draft.pages?.[pageIndex];
            const content = page?.dashboardView_Organizations?.content;

            const emptyRow = produce(content?.[0], (draft) => {
              if (!draft) return;
              draft.id = Math.random().toString();
              draft.name = '';
              draft.relationshipStages = [];
              draft.website = '';
              draft.owner = null;
              draft.lastTouchPointTimelineEvent = null;
              draft.accountDetails = null;
            });

            if (!emptyRow) return;
            content?.unshift(emptyRow);
          });
        },
      );

      setOrganizationsMeta((prev) =>
        produce(prev, (draft) => {
          draft.getOrganization.pagination.page = 1;
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        context?.previousEntries,
      );
    },
    onSuccess: ({ organization_Create: { id } }) => {
      push(`/organization/${id}`);
    },
    onSettled: () => {
      queryClient.invalidateQueries(queryKey);
    },
  });

  const hideOrganizations = useHideOrganizationsMutation(client, {
    onMutate: ({ ids }) => {
      const pageIndex = organizationsMeta.getOrganization.pagination.page - 1;
      queryClient.cancelQueries(queryKey);

      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            if (!draft) return;

            const page = draft.pages?.[pageIndex];
            const content = page?.dashboardView_Organizations?.content;
            const filteredContent = content?.filter(
              (o) => !(ids as string[]).includes(o.id),
            );

            set(
              draft,
              `pages[${pageIndex}].dashboardView_Organizations.content`,
              filteredContent,
            );
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
      queryClient.invalidateQueries(queryKey);
    },
  });

  const mergeOrganizations = useMergeOrganizationsMutation(client, {
    onMutate: ({ primaryOrganizationId, mergedOrganizationIds }) => {
      queryClient.cancelQueries(queryKey);

      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            if (!draft) return;

            const content =
              draft.pages?.[0]?.dashboardView_Organizations?.content;
            const targetOrganization = content?.find(
              (o) => o.id === primaryOrganizationId,
            );
            const filteredContent = [
              targetOrganization,
              ...(content ?? []).filter(
                (o) =>
                  ![
                    primaryOrganizationId,
                    ...(mergedOrganizationIds as string[]),
                  ].includes(o.id),
              ),
            ];

            if (!targetOrganization) return;

            set(
              draft,
              `pages[0].dashboardView_Organizations.content`,
              filteredContent,
            );
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
      queryClient.invalidateQueries(queryKey);
    },
  });

  return {
    createOrganization,
    hideOrganizations,
    mergeOrganizations,
  };
};
