'use client';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Image } from '@ui/media/Image';
import { VStack } from '@ui/layout/Stack';
import { GridItem } from '@ui/layout/Grid';

import { SidenavItem } from './SidenavItem';
import React, { useEffect } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { useLocalStorage } from 'usehooks-ts';

interface RootSidenavProps {
  isOwner: boolean;
}

export const RootSidenav = ({ isOwner }: RootSidenavProps) => {
  const router = useRouter();
  const pathname = usePathname();
  const [activeTab, setActiveTab] = useLocalStorage(
    `customeros-player-last-position`,
    { root: 'organization' },
  );

  useEffect(() => {
    if (pathname === '/') {
      setActiveTab({ ...activeTab, root: 'organization' });
    }
    if (pathname && pathname !== '/') {
      setActiveTab({ ...activeTab, root: pathname.substring(1) });
    }
  }, []);
  const handleItemClick = (path: string) => {
    setActiveTab({ ...activeTab, root: path });
    router.push(`/${path}`);
  };

  const checkIsActive = (path: string) => pathname?.startsWith(`/${path}`);

  return (
    <GridItem
      px='2'
      py='4'
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
        mb='4'
        tabIndex={0}
        role='button'
        cursor='pointer'
        justify='center'
        overflow='hidden'
        position='relative'
      >
        <Image
          width={40}
          height={40}
          w='40px'
          h='40px'
          alt='CustomerOS'
          pointerEvents='none'
          src='/logos/customer-os.png'
          transition='opacity 0.25s ease-in-out'
        />
      </Flex>

      <VStack spacing='2' w='full'>
        <SidenavItem
          label='Organizations'
          isActive={checkIsActive('organization')}
          onClick={() => handleItemClick('organization')}
          icon={(isActive) => (
            <Icons.Building7
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        <SidenavItem
          label='Customers'
          isActive={checkIsActive('customers')}
          onClick={() => handleItemClick('customers')}
          icon={(isActive) => (
            <Icons.CheckHeart
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        {isOwner && (
          <SidenavItem
            label='My portfolio'
            isActive={checkIsActive('portfolio')}
            onClick={() => handleItemClick('portfolio')}
            icon={(isActive) => (
              <Icons.Briefcase1
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
        <SidenavItem
          label='Settings'
          isActive={checkIsActive('settings')}
          onClick={() => router.push('/settings')}
          icon={(isActive) => (
            <Icons.Settings
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
      </VStack>
    </GridItem>
  );
};
