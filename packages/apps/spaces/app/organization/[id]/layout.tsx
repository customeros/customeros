import { Flex } from '@ui/layout/Flex';
import { GridItem } from '@ui/layout/Grid';
import { PageLayout } from '@shared/components/PageLayout';

import { OrganizationSidenav } from './src/components/OrganizationSidenav';

export default function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <PageLayout>
      <OrganizationSidenav />
      <GridItem h='100%' area='content' overflow='hidden' columnGap={2} gap={0}>
        <Flex flexDir='row' h='100%'>
          {children}
        </Flex>
      </GridItem>
    </PageLayout>
  );
}
