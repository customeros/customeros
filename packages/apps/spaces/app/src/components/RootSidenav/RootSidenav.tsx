'use client';
import { useEffect } from 'react';
import { produce } from 'immer';
import { signOut } from 'next-auth/react';
import { useLocalStorage } from 'usehooks-ts';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Image } from '@ui/media/Image';
import { VStack } from '@ui/layout/Stack';
import { GridItem } from '@ui/layout/Grid';

import { LogOut01 } from '@ui/media/icons/LogOut01';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import logo from './customer-os.png';
import { SidenavItem } from './components/SidenavItem';
import { GoogleSidebarNotification } from './components/GoogleSidebarNotification';

export const RootSidenav = () => {
  const client = getGraphQLClient();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [_, setOrganizationsMeta] = useOrganizationsMeta();
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { root: 'organization' },
  );

  const { data } = useGlobalCacheQuery(client);
  const globalCache = data?.global_Cache;

  const handleItemClick = (path: string) => {
    setLastActivePosition({ ...lastActivePosition, root: path });
    setOrganizationsMeta((prev) =>
      produce(prev, (draft) => {
        draft.getOrganization.pagination.page = 1;
      }),
    );

    router.push(`/${path}`);
  };

  const checkIsActive = (path: string, options?: { preset: string }) => {
    const [_pathName, _searchParams] = path.split('?');
    const presetParam = new URLSearchParams(searchParams?.toString()).get(
      'preset',
    );

    if (options?.preset) {
      return (
        pathname?.startsWith(`/${_pathName}`) && presetParam === options.preset
      );
    } else {
      return pathname?.startsWith(`/${_pathName}`) && !presetParam;
    }
  };

  const handleSignOutClick = () => {
    signOut();
  };

  useEffect(() => {
    const params = searchParams?.toString();

    if (pathname === '/') {
      setLastActivePosition({
        ...lastActivePosition,
        root: `organizations?${params}`,
      });
    }
    if (pathname && pathname !== '/') {
      setLastActivePosition({
        ...lastActivePosition,
        root: `${pathname.substring(1)}?${params}`,
      });
    }
  }, []);

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
          src={logo}
          transition='opacity 0.25s ease-in-out'
        />
      </Flex>

      <VStack spacing='2' w='full'>
        <SidenavItem
          label='Organizations'
          isActive={checkIsActive('organizations')}
          onClick={() => handleItemClick('organizations')}
          icon={(isActive) => (
            <Icons.Building7
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        <SidenavItem
          label='Customers'
          isActive={checkIsActive('organizations', { preset: 'customer' })}
          onClick={() => handleItemClick('organizations?preset=customer')}
          icon={(isActive) => (
            <Icons.CheckHeart
              boxSize='6'
              color={isActive ? 'gray.700' : 'gray.500'}
            />
          )}
        />
        {globalCache?.isOwner && (
          <SidenavItem
            label='My portfolio'
            isActive={checkIsActive('organizations', { preset: 'portfolio' })}
            onClick={() => handleItemClick('organizations?preset=portfolio')}
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
        <GoogleSidebarNotification />

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
        <SidenavItem
          label='Sign out'
          isActive={false}
          onClick={handleSignOutClick}
          icon={() => <LogOut01 boxSize='6' color='gray.500' />}
        />
      </VStack>
    </GridItem>
  );
};
