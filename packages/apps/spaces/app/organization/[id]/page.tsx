import { notFound } from 'next/navigation';

import { HydrationBoundary } from '@tanstack/react-query';

import { getDehydratedState } from '@shared/util/getDehydratedState';
import { TimelineContextsProvider } from '@organization/src/components/TimelineContextsProvider';

import { SideSection } from './src/components/SideSection';
import { MainSection } from './src/components/MainSection';
import { Panels, TabsContainer } from './src/components/Tabs';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from './src/graphql/organization.generated';
import { OrganizationTimelineWithActionsContext } from './src/components/Timeline/OrganizationTimelineWithActionsContext';

interface OrganizationPageProps {
  params: { id: string };
  searchParams: { tab?: string };
}

export default async function OrganizationPage({
  searchParams,
  params,
}: OrganizationPageProps) {
  const dehydratedClient = await getDehydratedState(useOrganizationQuery, {
    variables: { id: params.id },
    fetcher: useOrganizationQuery.fetcher,
  });

  const organizationData = (
    dehydratedClient.queries?.[0]?.state.data as OrganizationQuery
  )?.organization;

  if (organizationData?.hide) {
    notFound();
  }

  return (
    <HydrationBoundary state={dehydratedClient}>
      <TimelineContextsProvider id={params.id}>
        <SideSection>
          <TabsContainer>
            <Panels tab={searchParams.tab ?? 'about'} />
          </TabsContainer>
        </SideSection>

        <MainSection>
          <OrganizationTimelineWithActionsContext />
        </MainSection>
      </TimelineContextsProvider>
    </HydrationBoundary>
  );
}
