import { PageLayout } from '@shared/components/PageLayout';
import { RootSidenav } from '@shared/components/RootSidenav/RootSidenav';

export default function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <PageLayout>
      <RootSidenav />
      <div className='h-full overflow-hidden'>{children}</div>
    </PageLayout>
  );
}
