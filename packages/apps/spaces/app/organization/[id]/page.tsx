import { SideSection } from './src/components/SideSection';
import { MainSection } from './src/components/MainSection';
import { Panels, TabsContainer } from './src/components/Tabs';
import { OrganizationTimelineWithActionsContext } from './src/components/Timeline/OrganizationTimelineWithActionsContext';
import {
  GetCanAccessOrganizationDocument,
  GetCanAccessOrganizationQuery,
} from '@organization/src/graphql/getCanAccessOrganization.generated';
import { getServerGraphQLClient } from '@shared/util/getServerGraphQLClient';
import NotFound from './src/components/NotFound/NotFound';
import { TimelineContextsProvider } from '@organization/src/components/TimelineContextsProvider';

interface OrganizationPageProps {
  params: { id: string };
  searchParams: { tab?: string };
}

export default async function OrganizationPage({
  searchParams,
  params,
}: OrganizationPageProps) {
  const client = getServerGraphQLClient();

  try {
    await client.request<GetCanAccessOrganizationQuery>(
      GetCanAccessOrganizationDocument,
      {
        id: params.id,
      },
    );
  } catch (error) {
    return <NotFound />;
  }

  return (
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
  );
}
