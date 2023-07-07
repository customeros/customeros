
import { dehydrate, Hydrate } from '@tanstack/react-query';
import getQueryClient from '@shared/util/getQueryClient';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { gql } from 'graphql-request';

import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';
import {
  OrganizationInfo,
  OrganizationLogo,
  OrganizationTabs,
  OrganizationHeader,
} from './components/OrganizationInfo';
import { OrganizationTimeline } from './components/Timeline/OrganizationTimeline';

const ORGANIZATION_QUERY = gql`
  query Organization($id: ID!) {
    organization(id: $id) {
      id
      name
      description
      domains
      subIndustry
      website
      industry
      isPublic
      market
      employees
      socials {
        id
        platformName
        url
      }
      relationshipStages {
        relationship
        stage
      }
    }
  }
`;

interface OrganizationPageProps {
  params: { id: string };
}

export default async function OrganizationPage({
  params,
}: OrganizationPageProps) {
  // const queryClient = getQueryClient();
  // await queryClient.prefetchQuery(['organization', params.id]);
  // const dehydratedState = dehydrate(queryClient);

  const graphqlClient = getGraphQLClient();
  try {
    const req = await graphqlClient.request(ORGANIZATION_QUERY, {
      id: params.id,
    });
    console.log(req);
  } catch (e) {
    console.log(e);
  }

  return (
    <>
      <SideSection>
        <OrganizationInfo>
          <OrganizationHeader>
            <OrganizationLogo src='/logos/bigquery.svg' />
          </OrganizationHeader>

          <OrganizationTabs />
        </OrganizationInfo>
      </SideSection>
      <MainSection>
        <OrganizationTimeline id={params.id} />
      </MainSection>
    </>
  );
}
