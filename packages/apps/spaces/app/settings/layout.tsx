import { SettingsSidenav } from './src/components/SettingsSidenav';
import { GridItem } from '@ui/layout/Grid';

export default async function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <SettingsSidenav />
      <GridItem h='100%' area='content' overflowX='hidden' overflowY='auto'>
        {children}
      </GridItem>
    </>
  );
}
