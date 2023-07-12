import { Hydrate } from '@tanstack/react-query';
import { getDehydratedState } from '@shared/util/getDehydratedState';

import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';

import {
  Tabs,
  TabsHeader,
  AboutPanel,
  PeoplePanel,
  UpNextPanel,
  AccountPanel,
  SuccessPanel,
  TabsContainer,
  OrganizationLogo,
} from './components/Tabs';

import { useOrganizationQuery } from './graphql/organization.generated';

interface OrganizationPageProps {
  params: { id: string };
}

export default async function OrganizationPage({
  params,
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
          <TabsHeader>
            <OrganizationLogo src='/logos/bigquery.svg' />
          </TabsHeader>

          <Tabs>
            <UpNextPanel />
            <AccountPanel />
            <SuccessPanel />
            <PeoplePanel />
            <AboutPanel />
          </Tabs>
        </TabsContainer>
      </SideSection>

      <MainSection></MainSection>
    </Hydrate>
  );
}
