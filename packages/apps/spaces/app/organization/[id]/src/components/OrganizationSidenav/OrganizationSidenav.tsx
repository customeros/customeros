'use client';
import { useRouter, useParams, useSearchParams } from 'next/navigation';

import { useLocalStorage } from 'usehooks-ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { VStack } from '@ui/layout/Stack';
import { Link } from '@ui/navigation/Link';
import { GridItem } from '@ui/layout/Grid';
import { Text } from '@ui/typography/Text';
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';

export const OrganizationSidenav = () => {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();
  const isOnboardingFeatureOn = useFeatureIsOn('onboarding-status');

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [params?.id as string]: 'tab=about' },
  );

  const graphqlClient = getGraphQLClient();
  const { data } = useOrganizationQuery(graphqlClient, {
    id: params?.id as string,
  });
  const parentOrg = data?.organization?.subsidiaryOf?.[0]?.organization;

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const urlSearchParams = new URLSearchParams(searchParams?.toString());
    urlSearchParams.set('tab', tab);

    setLastActivePosition({
      ...lastActivePosition,
      [params?.id as string]: urlSearchParams.toString(),
    });
    router.push(`?${urlSearchParams}`);
  };

  return (
    <GridItem
      px='2'
      py='4'
      h='full'
      w='200px'
      display='flex'
      flexDir='column'
      gridArea='sidebar'
      background='white'
      position='relative'
      borderRight='1px solid'
      borderColor='gray.200'
    >
      <Flex gap='2' align='center' mb='4'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Go back'
          onClick={() => {
            router.push(`/${lastActivePosition?.root || 'organization'}`);
          }}
          icon={<Icons.ArrowNarrowLeft color='gray.700' boxSize='6' />}
        />

        <Flex flexDir='column'>
          {parentOrg && (
            <Link
              fontSize='xs'
              noOfLines={1}
              wordBreak='keep-all'
              href={`/organization/${parentOrg.id}?tab=about`}
            >
              {parentOrg.name}
            </Link>
          )}
          <Tooltip label={data?.organization?.name} placement='bottom'>
            <Text
              fontSize='lg'
              fontWeight='semibold'
              color='gray.700'
              noOfLines={1}
              wordBreak='keep-all'
            >
              {data?.organization?.name || 'Organization'}
            </Text>
          </Tooltip>
        </Flex>
      </Flex>

      <VStack spacing='1' w='full'>
        <SidenavItem
          label='About'
          isActive={checkIsActive('about') || !searchParams?.get('tab')}
          onClick={handleItemClick('about')}
          icon={
            <Icons.InfoSquare
              color={checkIsActive('about') ? 'gray.700' : 'gray.500'}
              boxSize='5'
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
              boxSize='5'
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
              boxSize='5'
            />
          }
        />
        <SidenavItem
          label='Issues'
          isActive={checkIsActive('issues')}
          onClick={handleItemClick('issues')}
          icon={
            <Ticket02
              color={checkIsActive('issues') ? 'gray.700' : 'gray.500'}
              boxSize='5'
            />
          }
        />
        {isOnboardingFeatureOn && (
          <SidenavItem
            label='Success'
            isActive={checkIsActive('success')}
            onClick={handleItemClick('success')}
            icon={
              <Trophy01
                color={checkIsActive('success') ? 'gray.700' : 'gray.500'}
                boxSize='5'
              />
            }
          />
        )}
      </VStack>
    </GridItem>
  );
};
