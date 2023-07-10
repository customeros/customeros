import { Hydrate } from '@tanstack/react-query';
import { getDehydratedState } from '@shared/util/getDehydratedState';

import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';
import {
  OrganizationLogo,
  OrganizationTabsHeader,
} from './components/OrganizationTabs/OrganizationTabsHeader';
import { OrganizationAboutTab } from './components/OrganizationAboutTab';
import { OrganizationUpNextTab } from './components/OrganizationUpNextTab';
import { OrganizationPeopleTab } from './components/OrganizationPeopleTab';
import { OrganizationAccountTab } from './components/OrganizationAccountTab';
import { OrganizationSuccessTab } from './components/OrganizationSuccessTab';
import { OrganizationTabs } from './components/OrganizationTabs/OrganizationTabs';
import { OrganizationTabsContainer } from './components/OrganizationTabs/OrganizationTabsContainer';

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
        <OrganizationTabsContainer>
          <OrganizationTabsHeader>
            <OrganizationLogo src='/logos/bigquery.svg' />
          </OrganizationTabsHeader>

          <OrganizationTabs>
            <OrganizationUpNextTab />
            <OrganizationAccountTab />
            <OrganizationSuccessTab />
            <OrganizationPeopleTab />
            <OrganizationAboutTab />
          </OrganizationTabs>
        </OrganizationTabsContainer>
      </SideSection>

      <MainSection></MainSection>
    </Hydrate>
  );
}
