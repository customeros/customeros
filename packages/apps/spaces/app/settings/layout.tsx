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
      <div className='h-full overflow-x-hidden overflow-y-auto'>{children}</div>
    </PageLayout>
  );
}
