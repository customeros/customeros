import { GridItem } from '@ui/layout/Grid';
import { PageLayout } from '@shared/components/PageLayout';

import { SettingsSidenav } from './src/components/SettingsSidenav';

export default async function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <PageLayout>
      <SettingsSidenav />
      <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
        {children}
      </GridItem>
    </PageLayout>
  );
}
