import { GridItem } from '@ui/layout/Grid';
import { PageLayout } from '@shared/components/PageLayout';
import { RootSidenav } from '@shared/components/RootSidenav/RootSidenav';

import { Providers } from './src/components/Providers';

export default function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Providers>
      <PageLayout>
        <RootSidenav />
        <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
          {children}
        </GridItem>
      </PageLayout>
    </Providers>
  );
}
