import { Flex } from '@ui/layout/Flex';
import { GridItem } from '@ui/layout/Grid';
import {TenantSidenav} from "./TenantSidenav/TenantSidenav";


export default async function OrganizationLayout({
  children,
}: {
  children: React.ReactNode;
}) {

  return (
    <>
      <TenantSidenav />
      <GridItem h='100%' area='content' overflow='hidden'>
        <Flex flexDir='row' gap='2'>
          {children}
        </Flex>
      </GridItem>
    </>
  );
}
