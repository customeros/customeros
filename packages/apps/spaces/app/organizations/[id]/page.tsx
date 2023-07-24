import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';

import { TabsContainer, Panels } from './components/Tabs';

import { OrganizationTimeline } from './components/Timeline/OrganizationTimeline';

interface OrganizationPageProps {
  params: { id: string };
  searchParams: { tab?: string };
}

export default async function OrganizationPage({
  searchParams,
}: OrganizationPageProps) {
  return (
    <>
      <SideSection>
        <TabsContainer>
          <Panels tab={searchParams.tab ?? 'about'} />
        </TabsContainer>
      </SideSection>

      <MainSection>
        <OrganizationTimeline />
      </MainSection>
    </>
  );
}
