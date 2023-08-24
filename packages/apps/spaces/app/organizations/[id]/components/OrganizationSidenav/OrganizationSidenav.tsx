'use client';
import { useRouter, useParams, useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Icons } from '@ui/media/Icon';
import { GridItem } from '@ui/layout/Grid';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { SidenavItem } from '@shared/components/RootSidenav/SidenavItem';
import { useOrganizationQuery } from '@organization/graphql/organization.generated';
import { useLocalStorage } from 'usehooks-ts';

export const OrganizationSidenav = () => {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [params?.id as string]: 'tab=about' },
  );

  const graphqlClient = getGraphQLClient();
  const { data } = useOrganizationQuery(graphqlClient, {
    id: params?.id as string,
  });

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const urlSearchParams = new URLSearchParams(searchParams ?? '');
    urlSearchParams.set('tab', tab);

    setLastActivePosition({
      ...lastActivePosition,
      [params?.id as string]: urlSearchParams.toString(),
    });
    // todo remove, for now needed
    router.push(`?${urlSearchParams}`);
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
      <Tooltip label={data?.organization?.name} placement='bottom'>
        <Flex gap='2' align='center' mb='4'>
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Go back'
            onClick={() => router.push(`/${lastActivePosition?.root || 'organization'}`)}
            icon={<Icons.ArrowNarrowLeft color='gray.700' boxSize='6' />}
          />

          <Text
            fontSize='lg'
            fontWeight='semibold'
            color='gray.700'
            noOfLines={1}
            wordBreak='keep-all'
          >
            {data?.organization?.name || 'Organization'}
          </Text>
        </Flex>
      </Tooltip>

      <VStack spacing='2' w='full'>
        <SidenavItem
          label='About'
          isActive={checkIsActive('about') || !searchParams?.get('tab')}
          onClick={handleItemClick('about')}
          icon={
            <Icons.InfoSquare
              color={checkIsActive('about') ? 'gray.700' : 'gray.500'}
              boxSize='6'
            />
          }
        />
        <SidenavItem
          label='People'
          isActive={checkIsActive('people')}
          onClick={handleItemClick('people')}
          icon={
            <Icons.Users2
              color={checkIsActive('people') ? 'gray.700' : 'gray.500'}
              boxSize='6'
            />
          }
        />
        <SidenavItem
          label='Account'
          isActive={checkIsActive('account')}
          onClick={handleItemClick('account')}
          icon={
            <Icons.ActivityHeart
              color={checkIsActive('account') ? 'gray.700' : 'gray.500'}
              boxSize='6'
            />
          }
        />
      </VStack>
    </GridItem>
  );
};
