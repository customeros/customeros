import { useRouter, useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Icons } from '@ui/media/Icon';
import { GridItem } from '@ui/layout/Grid';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';

import { SidenavItem } from './SidenavItem';

export const OrganizationSidenav = () => {
  const router = useRouter();
  const searchParams = useSearchParams();

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const params = new URLSearchParams(searchParams ?? '');
    params.set('tab', tab);

    router.push(`?${params}`);
  };

  return (
    <GridItem
      px='4'
      pt='4'
      pb='8'
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
      <Flex gap='2' align='center' mb='4'>
        <IconButton
          size='md'
          variant='ghost'
          aria-label='Go back'
          onClick={() => router.push('/organization')}
          icon={<Icons.ArrowNarrowLeft color='gray.700' boxSize='6' />}
        />

        <Text
          fontSize='lg'
          fontWeight='bold'
          color='gray.700'
          noOfLines={1}
          wordBreak='keep-all'
        >
          Organization Lorem ipsum sin dolor
        </Text>
      </Flex>

      <VStack spacing='1' w='full'>
        <SidenavItem
          label='About'
          isActive={checkIsActive('about')}
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
            <Icons.Folder
              color={checkIsActive('account') ? 'gray.700' : 'gray.500'}
              boxSize='6'
            />
          }
        />
      </VStack>
    </GridItem>
  );
};
