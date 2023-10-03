import set from 'lodash/set';
import { produce } from 'immer';
import { useRouter } from 'next/navigation';
import { RowSelectionState } from '@tanstack/react-table';
import { useQueryClient, InfiniteData } from '@tanstack/react-query';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  GetOrganizationsQuery,
  useInfiniteGetOrganizationsQuery,
} from '../graphql/getOrganizations.generated';
import { useOrganizationsMeta } from '../shared/state';
import { useCreateOrganizationMutation } from '../graphql/createOrganization.generated';
import { useHideOrganizationsMutation } from '../graphql/hideOrganizations.generated';
import { useMergeOrganizationsMutation } from '../graphql/mergeOrganizations.generated';

interface UseOrganizationsPageMethodsOptions {
  selection: RowSelectionState;
  setEnableSelection: React.Dispatch<React.SetStateAction<boolean>>;
}

export const useOrganizationsPageMethods = ({
  selection,
  setEnableSelection,
}: UseOrganizationsPageMethodsOptions) => {
  const { push } = useRouter();
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();

  const queryKey = useInfiniteGetOrganizationsQuery.getKey(
    organizationsMeta.getOrganization,
  );

  const createOrganization = useCreateOrganizationMutation(client, {
    onSuccess: ({ organization_Create: { id } }) => {
      push(`/organizations/${id}`);
    },
  });

  const hideOrganizations = useHideOrganizationsMutation(client, {
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

            const keys = Object.keys(selection);
            const page = draft.pages?.[pageIndex];
            const content = page?.dashboardView_Organizations?.content;
            const filteredContent = content?.filter(
              (_, idx) => !keys.includes(String(idx)),
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
    onMutate: () => {
      setEnableSelection(false);
      queryClient.cancelQueries(queryKey);

      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsQuery>>(queryKey);

      queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            if (!draft) return;

            const [target, ...rest] = Object.keys(selection);
            const page = draft.pages?.[0];
            const content = page?.dashboardView_Organizations?.content;
            const targetOrganization = content?.[Number(target)];
            const filteredContent = content?.filter(
              (_, idx) => !rest?.includes(String(idx)),
            );

            if (!targetOrganization) return;
            filteredContent?.splice(Number(target), 1);
            filteredContent?.unshift(targetOrganization);

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
