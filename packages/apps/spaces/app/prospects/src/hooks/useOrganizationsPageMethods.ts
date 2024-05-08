import { produce } from 'immer';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
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

  const updateOrganization = useUpdateOrganizationMutation(client, {
    onMutate: (payload) => {
      queryClient.cancelQueries({ queryKey });

      const previousOrganizations =
        queryClient.getQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
          queryKey,
        );

      queryClient.setQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
        queryKey,
        (old) => {
          const pageIndex = 0;

          return produce(old, (draft) => {
            const content =
              draft?.pages?.[pageIndex]?.dashboardView_Organizations?.content;
            const index = content?.findIndex(
              (item) => item.metadata.id === payload.input.id,
            );
            console.log('🏷️ ----- payload: ', payload);

            if (content && index !== undefined && index > -1) {
              content[index].stage = payload.input.stage;
            }
          });
        },
      );

      return { previousOrganizations };
    },
    onError: (_, __, context) => {
      if (context?.previousOrganizations) {
        queryClient.setQueryData<InfiniteData<GetOrganizationsKanbanQuery>>(
          queryKey,
          context.previousOrganizations,
        );
      }
    },
    onSettled: () =>
      queryClient.invalidateQueries({
        queryKey,
      }),
  });

  return {
    createOrganization,
    updateOrganization,
  };
};
