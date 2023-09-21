import { GridItem } from '@ui/layout/Grid';
import { SettingsSidenav } from './SettingsSidenav';

export default async function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <SettingsSidenav />
      <GridItem h='100%' area='content' overflow='hidden'>
        {children}
      </GridItem>
    </>
  );
}
