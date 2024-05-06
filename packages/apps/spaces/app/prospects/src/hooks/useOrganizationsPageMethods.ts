import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useCreateOrganizationMutation } from '@organizations/graphql/createOrganization.generated';

import { GetOrganizationsKanbanQuery } from '../graphql/getOrganizationsKanban.generated';

export const useOrganizationsPageMethods = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const queryKey = ['getOrganizationKanban.infinite'];
  const createOrganization = useCreateOrganizationMutation(client, {
    onMutate: () => {
      const pageIndex = 0;
      queryClient.cancelQueries({ queryKey });

      const previousEntries =
        queryClient.getQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
          queryKey,
        );

      queryClient.setQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
        queryKey,
        (old) => {
          return produce(old, (draft) => {
            if (!draft) return;

            const page = draft.pages?.[pageIndex];
            let content = page?.dashboardView_Organizations?.content;

            const emptyRow = produce(content?.[0], (draft) => {
              if (!draft) return;
              draft.metadata.id = Math.random().toString();
              draft.name = '';
              draft.owner = null;
              draft.accountDetails = null;
            });

            if (!emptyRow) return;
            content = [emptyRow, ...(content ?? [])];
          });
        },
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
        queryKey,
        context?.previousEntries,
      );
      toastError(
        `We couldn't create this organization`,
        'create-organization-error',
      );
    },
    onSuccess: ({ organization_Create: { id } }) => {
      toastSuccess(`Organization created`, 'create-organization-success');
      queryClient.invalidateQueries({ queryKey });
    },
  });

  return {
    createOrganization,
  };
};
