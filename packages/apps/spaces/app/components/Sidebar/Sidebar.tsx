'use client';
import { signOut } from 'next-auth/react';
import { useJune } from '@spaces/hooks/useJune';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Image } from '@ui/media/Image';
import { VStack } from '@ui/layout/Stack';
import { GridItem } from '@ui/layout/Grid';

import { SidebarItem } from './SidebarItem';

interface SidebarProps {
  isOwner: boolean;
}

export const Sidebar = ({ isOwner }: SidebarProps) => {
  const analytics = useJune();

  const handleClick = () => {
    analytics?.reset();
    signOut();
  };

  return (
    <GridItem
      px='4'
      py='8'
      h='full'
      w='200px'
      bg='white'
      display='flex'
      flexDir='column'
      gridArea='sidebar'
      position='relative'
      border='1px solid'
      borderRadius='2xl'
      borderColor='gray.200'
    >
      <Flex
        mb='2'
        tabIndex={0}
        role='button'
        cursor='pointer'
        justify='center'
        overflow='hidden'
        position='relative'
      >
        <Image
          width={32}
          height={32}
          w='32px'
          h='32px'
          alt='Openline'
          pointerEvents='none'
          src='/logos/openline_small.svg'
          transition='opacity 0.25s ease-in-out'
        />
      </Flex>

      <VStack spacing='1' w='full'>
        <SidebarItem
          href='/organization'
          label='Organizations'
          icon={(isActive) => (
            <Icons.Building7
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        <SidebarItem
          href='/customers'
          label='Customers'
          icon={(isActive) => (
            <Icons.CheckHeart
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        {isOwner && (
          <SidebarItem
            href='/portfolio'
            label='My portfolio'
            icon={(isActive) => (
              <Icons.ClipboardCheck
                boxSize='6'
                color={isActive ? 'gray.700' : 'gray.500'}
              />
            )}
          />
        )}
      </VStack>

      <VStack
        spacing='1'
        flexDir='column'
        flexWrap='initial'
        flexGrow='1'
        justifyContent='flex-end'
      >
        <SidebarItem
          href='/settings'
          label='Settings'
          icon={(isActive) => (
            <Icons.Settings
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        <SidebarItem
          label='Logout'
          onClick={handleClick}
          icon={(isActive) => (
            <Icons.Logout1
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
      </VStack>
    </GridItem>
  );
};
