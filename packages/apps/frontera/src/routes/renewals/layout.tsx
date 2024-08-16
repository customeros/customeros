import { PageLayout } from '@shared/components/PageLayout';
import { RootSidenav } from '@shared/components/RootSidenav/RootSidenav';

export default function RenewalsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <PageLayout isResizable>
      <RootSidenav />
      <div className='h-full overflow-hidden'>{children}</div>
    </PageLayout>
  );
}
