'use client';
import { useRouter } from 'next/navigation';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { SearchSm } from '@ui/media/icons/SearchSm';

import HalfCirclePattern from './src/assets/HalfCirclePattern';

export default function NotFound() {
  const router = useRouter();
  return (
    <Box
      p={0}
      flex={1}
      as={Flex}
      flexDirection='column'
      bgImage='/backgrounds/organization/half-circle-pattern.svg'
      backgroundRepeat='no-repeat'
      backgroundSize='contain'
      h='100vh'
      w='100vw'
      position='relative'
      alignItems='center'
      justifyContent='center'
    >
      <Box
        position='absolute'
        height='50vh'
        maxH='768px'
        width='768px'
        top='50%'
        left='50%'
        style={{
          transform: 'translate(-50%, -90%) rotate(180deg)',
        }}
      >
        <HalfCirclePattern />
      </Box>
      <Flex
        position='relative'
        direction='column'
        alignItems='center'
        justifyContent='center'
        h='50vh'
      >
        <FeaturedIcon colorScheme='primary' size='lg'>
          <Box>
            <SearchSm />
          </Box>
        </FeaturedIcon>
        <Heading fontWeight={600} fontSize='6xl' color='gray.900' py={6}>
          We lost this page
        </Heading>
        <Text color='gray.600' fontSize='xl' pb={12} px={8}>
          There was a small hiccup in the success plan. Let’s get you back to a
          familiar place.
        </Text>
        <Button
          colorScheme='primary'
          variant='outline'
          size='lg'
          onClick={() => router.push('/')}
        >
          Take me home
        </Button>
      </Flex>
    </Box>
  );
}
