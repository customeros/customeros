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
      <div className='h-full overflow-hidden'>
        <div className='flex h-full'>{children}</div>
      </div>
    </PageLayout>
  );
}
