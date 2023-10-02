import { PageLayout } from '@shared/components/PageLayout';

export default async function OrganizationsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <PageLayout>{children}</PageLayout>;
}
