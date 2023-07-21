import { Hydrate } from '@tanstack/react-query';
import { getDehydratedState } from '@shared/util/getDehydratedState';

import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';

import { TabsContainer, Panels } from './components/Tabs';

import { useOrganizationQuery } from './graphql/organization.generated';
import {OrganizationTimeline} from "./components/Timeline/OrganizationTimeline";

interface OrganizationPageProps {
  params: { id: string };
  searchParams: { tab?: string };
}

export default async function OrganizationPage({
  params,
  searchParams,
}: OrganizationPageProps) {
  const variables = { id: params.id };
  const dehydratedState = await getDehydratedState(
    useOrganizationQuery,
    variables,
  );

  return (
    <Hydrate state={dehydratedState}>
      <SideSection>
        <TabsContainer>
          <Panels tab={searchParams.tab ?? 'about'} />
        </TabsContainer>
      </SideSection>

      <MainSection>

        <OrganizationTimeline />
      </MainSection>
    </Hydrate>
  );
}
