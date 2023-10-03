import { RootSidenav } from 'app/components/RootSidenav/RootSidenav';

import { GridItem } from '@ui/layout/Grid';

export default function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <RootSidenav />
      <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
        {children}
      </GridItem>
    </>
  );
}
