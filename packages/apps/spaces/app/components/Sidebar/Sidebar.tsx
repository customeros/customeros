import Company from '@spaces/atoms/icons/Company';
import Settings from '@spaces/atoms/icons/Settings';
import Contacts from '@spaces/atoms/icons/Contacts';
import Customer from '@spaces/atoms/icons/Customer';
import Portfolio from '@spaces/atoms/icons/Portfolio';

import { SidebarItem } from './SidebarItem';
import { LogoutSidebarItem } from './LogoutSidebarItem';

import { GridItem } from '@ui/layout/Grid';
import { Flex } from '@ui/layout/Flex';
import { Image } from '@ui/media/Image';

interface SidebarProps {
  isOwner: boolean;
}

export const Sidebar = ({ isOwner }: SidebarProps) => {
  return (
    <GridItem
      h='full'
      w='80px'
      shadow='base'
      bg='#f0f0f0'
      display='flex'
      flexDir='column'
      gridArea='sidebar'
      position='relative'
    >
      <Flex
        mb='4'
        pb='4'
        pt='8'
        tabIndex={0}
        role='button'
        cursor='pointer'
        justify='center'
        overflow='hidden'
        position='relative'
      >
        <Image
          width={31}
          height={40}
          w='31px'
          h='40px'
          alt='Openline'
          pointerEvents='none'
          src='/logos/openline_small.svg'
          transition='opacity 0.25s ease-in-out'
        />
      </Flex>

      <SidebarItem
        href='/organization'
        label='Organizations'
        icon={<Company height={24} width={24} style={{ scale: '0.8' }} />}
      />
      <SidebarItem
        href='/customers'
        label='Customers'
        icon={<Customer height={24} width={24} style={{ scale: '0.8' }} />}
      />
      {isOwner && (
        <SidebarItem
          href='/portfolio'
          label='My portfolio'
          icon={<Portfolio height={24} width={24} style={{ scale: '0.8' }} />}
        />
      )}

      <Flex
        mb='4'
        flexDir='column'
        flexWrap='initial'
        flexGrow='1'
        justifyContent='flex-end'
      >
        <SidebarItem
          href='/settings'
          label='Settings'
          icon={<Settings height={24} width={24} style={{ scale: '0.8' }} />}
        />
        <LogoutSidebarItem />
      </Flex>
    </GridItem>
  );
};
