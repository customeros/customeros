import { Hydrate } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { GridItem } from '@ui/layout/Grid';
import { getDehydratedState } from '@shared/util/getDehydratedState';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

import { OrganizationSidenav } from './components/OrganizationSidenav';

export default async function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const dehydratedState = await getDehydratedState(useTenantNameQuery);

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
