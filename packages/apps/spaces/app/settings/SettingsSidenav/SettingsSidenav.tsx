'use client';
import { useRouter, useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Icons } from '@ui/media/Icon';
import { GridItem } from '@ui/layout/Grid';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';

import { SidenavItem } from '@shared/components/RootSidenav/SidenavItem';
import { useLocalStorage } from 'usehooks-ts';

export const SettingsSidenav = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { ['settings']: 'oauth', root: 'organization' },
  );

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const params = new URLSearchParams(searchParams ?? '');
    params.set('tab', tab);
    setLastActivePosition({ ...lastActivePosition, settings: tab });
    // todo remove, for now needed
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
      border='1px solid'
      borderRadius='2xl'
      borderColor='gray.200'
    >
      <Flex gap='2' align='center' mb='4'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Go back'
          onClick={() => router.push(`/${lastActivePosition.root}`)}
          icon={<Icons.ArrowNarrowLeft color='gray.700' boxSize='6' />}
        />

        <Text
          fontSize='lg'
          fontWeight='semibold'
          color='gray.700'
          noOfLines={1}
          wordBreak='keep-all'
        >
          Tenant Settings
        </Text>
      </Flex>

      <VStack spacing='2' w='full'>
        <SidenavItem
          label='OAuth Settings'
          isActive={checkIsActive('oauth') || !searchParams?.get('tab')}
          onClick={handleItemClick('oauth')}
          icon={
            <Icons.InfoSquare
              color={checkIsActive('oauth') ? 'gray.700' : 'gray.500'}
              boxSize='6'
            />
          }
        />
        <SidenavItem
          label='Billing'
          isActive={checkIsActive('billing')}
          onClick={handleItemClick('billing')}
          icon={
            <Icons.Folder
              color={checkIsActive('billing') ? 'gray.700' : 'gray.500'}
              boxSize='6'
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
              boxSize='6'
            />
          }
        />
      </VStack>
    </GridItem>
  );
};
