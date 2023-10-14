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
import { GraphQLClient } from 'graphql-request';

interface OrganizationPageProps {
  params: { id: string };
  searchParams: { tab?: string };
}

export default async function OrganizationPage({
  searchParams,
  params,
}: OrganizationPageProps) {
  const client = getServerGraphQLClient();
  const result = await fetchOrganizationAccess({ client, id: params.id });

  if (!result) {
    return <NotFound />;
  }

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

interface FetchOrganizationAccessArgs {
  client: GraphQLClient;
  id: string;
}

async function fetchOrganizationAccess({
  client,
  id,
}: FetchOrganizationAccessArgs): Promise<GetCanAccessOrganizationQuery | null> {
  try {
    return await client.request<GetCanAccessOrganizationQuery>(
      GetCanAccessOrganizationDocument,
      {
        id,
      },
    );
  } catch (error) {
    return null;
  }
}
