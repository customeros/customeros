import { TimelineContextsProvider } from '@organization/src/components/TimelineContextsProvider';

import { SideSection } from './src/components/SideSection';
import { MainSection } from './src/components/MainSection';
import { Panels, TabsContainer } from './src/components/Tabs';
import { OrganizationTimelineWithActionsContext } from './src/components/Timeline/OrganizationTimelineWithActionsContext';

interface OrganizationPageProps {
  params: { id: string };
  searchParams: { tab?: string };
}

export default async function OrganizationPage({
  searchParams,
  params,
}: OrganizationPageProps) {
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
