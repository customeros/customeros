import { Flex } from '@ui/layout/Flex';
import { GridItem } from '@ui/layout/Grid';
import {SettingsSidenav} from "./SettingsSidenav/SettingsSidenav";


export default async function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {

  return (
    <>
      <SettingsSidenav />
      <GridItem h='100%' area='content' overflow='hidden'>
        <Flex flexDir='row' gap='2'>
          {children}
        </Flex>
      </GridItem>
    </>
  );
}
