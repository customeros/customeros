import { SideSection } from './src/components/SideSection';
import { MainSection } from './src/components/MainSection';
import { TabsContainer, Panels } from './src/components/Tabs';
import { OrganizationTimelineWithActionsContext } from '@organization/src/components/Timeline/OrganizationTimelineWithActionsContext';

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
        <OrganizationTimelineWithActionsContext />
      </MainSection>
    </>
  );
}
