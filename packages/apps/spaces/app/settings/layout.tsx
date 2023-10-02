import { SettingsSidenav } from './SettingsSidenav';
import { PageLayout } from '@shared/components/PageLayout';

export default async function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <PageLayout>
      <SettingsSidenav />
      {children}
    </PageLayout>
  );
}
