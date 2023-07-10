import { Hydrate } from '@tanstack/react-query';
import { getDehydratedState } from '@shared/util/getDehydratedState';

import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';
import {
  OrganizationInfo,
  OrganizationLogo,
  OrganizationTabs,
  OrganizationHeader,
} from './components/OrganizationInfo';
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
        <OrganizationInfo>
          <OrganizationHeader>
            <OrganizationLogo src='/logos/bigquery.svg' />
          </OrganizationHeader>

          <OrganizationTabs />
        </OrganizationInfo>
      </SideSection>
      <MainSection>{params.id}</MainSection>
    </Hydrate>
  );
}
