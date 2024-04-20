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
        <div className='h-full overflow-hidden'>{children}</div>
      </PageLayout>
    </Providers>
  );
}
