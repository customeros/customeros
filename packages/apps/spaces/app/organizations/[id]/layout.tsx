import { Hydrate } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { GridItem } from '@ui/layout/Grid';
import { getDehydratedState } from '@shared/util/getDehydratedState';
import { useOrganizationQuery } from '@organization/graphql/organization.generated';

import { OrganizationSidenav } from './components/OrganizationSidenav';

export default async function OrganizationLayout({
  children,
  params: { id },
}: {
  children: React.ReactNode;
  params: { id: string };
}) {
  const dehydratedState = await getDehydratedState(useOrganizationQuery, {
    id,
  });

  return (
    <Hydrate state={dehydratedState}>
      <OrganizationSidenav />
      <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
        <Flex flexDir='row' gap='2'>
          {children}
        </Flex>
      </GridItem>
    </Hydrate>
  );
}
