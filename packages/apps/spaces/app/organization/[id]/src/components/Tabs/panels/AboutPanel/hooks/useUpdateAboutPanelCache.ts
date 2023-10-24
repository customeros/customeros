import { InfiniteData, QueryKey, useQueryClient } from '@tanstack/react-query';
import { GetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { GetOrganizationsQuery } from '@organizations/graphql/getOrganizations.generated';
import { OrganizationQuery } from '@organization/src/graphql/organization.generated';
import { Organization } from '@graphql/types';

export function useUpdateAboutPanelCache() {
  const queryClient = useQueryClient();

  return async (updatedData: Organization, queryKey: QueryKey) => {
    await queryClient.cancelQueries({ queryKey });

    queryClient.setQueryData<OrganizationQuery>(
      queryKey,
      (currentCache): OrganizationQuery => {
        return {
          ...currentCache,
          organization: {
            ...currentCache?.organization,
            ...updatedData,
          },
        };
      },
    );

    return;
  };
}

export function useUpdateOrganizationInTableCache() {
  const queryClient = useQueryClient();
  const [organizationsMeta] = useOrganizationsMeta();
  const queryKey = [
    'getOrganizations.infinite',
    organizationsMeta.getOrganization,
  ];

  return async (updatedEvent: Organization) => {
    await queryClient.cancelQueries({ queryKey });
    queryClient.setQueryData<InfiniteData<GetOrganizationsQuery>>(
      queryKey,
      (currentCache): InfiniteData<GetOrganizationsQuery> => {
        const updatedPages = currentCache?.pages?.map((page) => {
          const updatedOrganizations =
            page?.dashboardView_Organizations?.content?.map(
              (event: Record<string, any>) => {
                if (event.id === updatedEvent?.id) {
                  return {
                    ...event,
                    ...updatedEvent,
                  };
                }
                return event;
              },
            );

          return {
            ...page,
            dashboardView_Organizations: {
              ...page.dashboardView_Organizations,
              content: [...(updatedOrganizations ?? [])],
            },
          };
        });
        return {
          ...currentCache,
          pages: updatedPages,
        } as InfiniteData<GetTimelineQuery>;
      },
    );
  };
}
