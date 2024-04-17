'use client';
import React from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

import { useLocalStorage } from 'usehooks-ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { VStack } from '@ui/layout/Stack';
import { GridItem } from '@ui/layout/Grid';
import { Text } from '@ui/typography/Text';
import { Map01 } from '@ui/media/icons/Map01';
import { IconButton } from '@ui/form/IconButton';
import { Receipt } from '@ui/media/icons/Receipt';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';

export const SettingsSidenav = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isMasterPlansEnabled = useFeatureIsOn('settings-master-plans-view');
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { ['settings']: 'oauth', root: 'organization' },
  );

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('tab', tab);
    setLastActivePosition({ ...lastActivePosition, settings: tab });
    router.push(`?${params}`);
  };

  return (
    <GridItem
      px='2'
      py='4'
      h='full'
      w='200px'
      background='gray.25'
      display='flex'
      flexDir='column'
      gridArea='sidebar'
      position='relative'
      borderRight='1px solid'
      borderColor='gray.200'
    >
      <Flex gap='2' align='center' mb='4'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Go back'
          onClick={() => router.push(`/${lastActivePosition.root}`)}
          icon={<Icons.ArrowNarrowLeft color='gray.700' boxSize='5' />}
        />

        <Text
          fontSize='lg'
          fontWeight='semibold'
          color='gray.700'
          noOfLines={1}
          wordBreak='keep-all'
        >
          Settings
        </Text>
      </Flex>

      <VStack spacing='2' w='full'>
        <SidenavItem
          label='My Account'
          isActive={checkIsActive('oauth') || !searchParams?.get('tab')}
          onClick={handleItemClick('oauth')}
          icon={
            <Icons.InfoSquare
              color={checkIsActive('oauth') ? 'gray.700' : 'gray.500'}
              boxSize='5'
            />
          }
        />
        <SidenavItem
          label='Customer billing'
          isActive={checkIsActive('billing')}
          onClick={handleItemClick('billing')}
          icon={
            <Receipt
              color={checkIsActive('billing') ? 'gray.700' : 'gray.500'}
              boxSize='5'
            />
          }
        />
        <SidenavItem
          label='Integrations'
          isActive={checkIsActive('integrations')}
          onClick={handleItemClick('integrations')}
          icon={
            <Icons.DataFlow3
              color={checkIsActive('integrations') ? 'gray.700' : 'gray.500'}
              boxSize='5'
            />
          }
        />
        {isMasterPlansEnabled && (
          <SidenavItem
            label='Master plans'
            isActive={checkIsActive('master-plans')}
            onClick={handleItemClick('master-plans')}
            icon={
              <Map01
                color={checkIsActive('master-plans') ? 'gray.700' : 'gray.500'}
                boxSize='5'
              />
            }
          />
        )}
      </VStack>
      <VStack
        spacing='1'
        flexDir='column'
        flexWrap='initial'
        flexGrow='1'
        justifyContent='flex-end'
        sx={{
          '& > span': {
            width: '100%',
          },
        }}
      >
        <NotificationCenter />
      </VStack>
      <Flex h='64px' />
    </GridItem>
  );
};
