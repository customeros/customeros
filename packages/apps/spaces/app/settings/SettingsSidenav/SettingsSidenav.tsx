'use client';
import { useRouter, useSearchParams, useParams } from 'next/navigation';

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

export const TenantSidenav = () => {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();
  const graphqlClient = getGraphQLClient();
  const { data } = useOrganizationQuery(graphqlClient, {
    id: params?.id as string,
  });

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const params = new URLSearchParams(searchParams ?? '');
    params.set('tab', tab);
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
      <Tooltip label={data?.organization?.name} placement='bottom'>
        <Flex gap='2' align='center' mb='4'>
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Go back'
            onClick={() => router.push('/organization')}
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
      </Tooltip>

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
          label='Billing Info'
          isActive={checkIsActive('billing')}
          onClick={handleItemClick('billing')}
          icon={
            <Icons.Folder
              color={checkIsActive('billing') ? 'gray.700' : 'gray.500'}
              boxSize='6'
            />
          }
        />
      </VStack>
    </GridItem>
  );
};
